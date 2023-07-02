import { useEffect, useState } from "react";
import { SessionOverview } from "./data";
import { loadSessionOverview } from "./dataloading";
import { MainApp } from "./mainapp";

export function App() {
  const [sessionOverview, setSessionOverview] = useState<SessionOverview | null>(null);

  const onSessionChange = () => {
    (async () => {
      const loadedSessionOverview = await loadSessionOverview();
      console.log(`Loaded the session overview`);
      setSessionOverview(loadedSessionOverview);
    })();
  };
  useEffect(() => {
    if (sessionOverview == null) {
      onSessionChange();
    }
  });

  if (sessionOverview == null) {
    return <div>Loading...</div>;
  }

  return <MainApp sessionOverview={sessionOverview} onSessionChange={onSessionChange} />;
}
