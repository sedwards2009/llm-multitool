import { useEffect, useState } from "react";
import { SessionOverview } from "./data";
import { loadSessionOverview } from "./dataloading";
import { MainApp } from "./mainapp";

export function App() {

  const [sessionOverview, setSessionOverview] = useState<SessionOverview | null>(null);

  useEffect(() => {
    if (sessionOverview == null) {
      (async () => {
        const loadedSessionOverview = await loadSessionOverview();
        console.log(`Loaded the session overview`);
        setSessionOverview(loadedSessionOverview);
      })();
    }
  });

  if (sessionOverview == null) {
    return <div>Loading...</div>;
  }

  return <MainApp sessionOverview={sessionOverview} />;
}
