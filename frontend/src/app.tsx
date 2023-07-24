import { useEffect, useState } from "react";
import { ModelOverview, SessionOverview } from "./data";
import { loadModelOverview, loadSessionOverview } from "./dataloading";
import { MainApp } from "./mainapp";

export function App() {
  const [sessionOverview, setSessionOverview] = useState<SessionOverview | null>(null);
  const [modelOverview, setModelOverview] = useState<ModelOverview | null>(null);

  const loadModelOverviewData = async () => {
    const overview = await loadModelOverview();
    console.log(`Loaded the model overview`);
    setModelOverview(overview);
  };

  const onSessionChange = () => {
    (async () => {
      const loadedSessionOverview = await loadSessionOverview();
      console.log(`Loaded the session overview`);
      setSessionOverview(loadedSessionOverview);
    })();
  };

  useEffect(() => {
    loadModelOverviewData();
    if (sessionOverview == null) {
      onSessionChange();
    }
  }, []);

  if (sessionOverview == null || modelOverview == null) {
    return <div>Loading...</div>;
  }

  return <MainApp
    modelOverview={modelOverview}
    sessionOverview={sessionOverview}
    onSessionChange={onSessionChange}
    />;
}
