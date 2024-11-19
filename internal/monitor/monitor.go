package monitor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/EtincelleCoworking/availability-monitor/pkg/sensor"
	"github.com/EtincelleCoworking/availability-monitor/pkg/ws"
)

var upgrader = websocket.Upgrader{}

type SensorManager interface {
	List() []sensor.Sensor
}
type MonitorService struct {
	hub                *ws.Hub
	sensorManager      SensorManager
	logger             *log.Logger
	unreachableTimeout time.Duration
	maybeFreeTimeout   time.Duration
	freeTimeout        time.Duration
}

func NewMonitorService(sensorManager SensorManager, logger *log.Logger,
	unreachableTimeout time.Duration,
	maybeFreeTimeout time.Duration,
	freeTimeout time.Duration) *MonitorService {
	hub := ws.NewHub()
	return &MonitorService{
		hub:                hub,
		sensorManager:      sensorManager,
		logger:             logger,
		unreachableTimeout: unreachableTimeout,
		maybeFreeTimeout:   maybeFreeTimeout,
		freeTimeout:        freeTimeout,
	}
}

func (n *MonitorService) MapHandlers(e *echo.Echo) {
	e.GET("/ws/", n.serveWs)
}

func (n *MonitorService) Run(ctx context.Context) {
	go n.hub.Run(ctx)
	go n.monitorStateChange(ctx)
}

type wsMessage struct {
	Location string `json:"location"`
	State    string `json:"state"`
}

func (n *MonitorService) monitorStateChange(ctx context.Context) {
	states := map[string]string{}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sensors := n.sensorManager.List()
			for _, s := range sensors {
				prevState := states[s.MAC()]
				if prevState != n.getState(s) {
					msg := wsMessage{
						Location: s.Location(),
						State:    n.getState(s),
					}
					j, err := json.Marshal(msg)
					if err != nil {
						n.logger.Println(err)
					} else {
						n.hub.Broadcast(j)
					}
				}
				states[s.MAC()] = n.getState(s)
			}
		}
	}

}

// serveWs handles websocket requests from the peer.
func (n *MonitorService) serveWs(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	client := ws.NewClient(n.hub, conn)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()

	n.sendSensorList(conn)
	return nil
}

func (n *MonitorService) sendSensorList(conn *websocket.Conn) {

	sensors := n.sensorManager.List()
	list := []wsMessage{}
	for _, s := range sensors {
		list = append(list, wsMessage{
			Location: s.Location(),
			State:    n.getState(s),
		})
	}

	j, _ := json.Marshal(list)
	n.logger.Printf("%s", j)

	conn.WriteMessage(1, j)
}
