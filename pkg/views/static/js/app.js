const states = {
    free: "🟢",
    maybe: "🟡",
    busy: "🔴",
    unreachable: "⚠️"
  }
  const ws = new WebSocket("/ws");
  const App = () => {
    const [sensors, setSensors] = React.useState([]);
    React.useEffect(() => {
      ws.onmessage = (ev) => {
        const msg = JSON.parse(ev.data);
        if(Array.isArray(msg)){
            setSensors(msg)
            return;
        }
        setSensors((sensors) => {
          const newSensors = [...sensors];
          let found = false;
          for (const sensor of newSensors) {
            if (sensor.location !== msg.location) continue;
            sensor.state = msg.state;
            found = true;
            break;
          }
          if (!found) {
            newSensors.push(msg);
          }
          return newSensors;
        });
      };
    });

    return (
      <div>
        <h1>Sensors ({sensors.length})</h1>
        <ul>
          {sensors.map(({ location, state }) => (
            <li key={location}>
              {location}: {state} {states[state]}
            </li>
          ))}
        </ul>
      </div>
    );
  };

  const root = ReactDOM.createRoot(document.getElementById("root"));
  root.render(<App />);