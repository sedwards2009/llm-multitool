import { useEffect, useState } from "react";
import { ModelOverview, SessionOverview, TemplateOverview } from "./data";
import { loadModelOverview, loadSessionOverview, loadTemplateOverview } from "./dataloading";
import { MainApp } from "./mainapp";

export function App() {
  const [sessionOverview, setSessionOverview] = useState<SessionOverview | null>(null);
  const [modelOverview, setModelOverview] = useState<ModelOverview | null>(null);
  const [templateOverview, setTemplateOverview] = useState<TemplateOverview | null>(null);

  const loadModelOverviewData = async () => {
    const overview = await loadModelOverview();
    console.log(`Loaded the model overview`);
    setModelOverview(overview);
  };

  const loadTemplateOverviewData = async () => {
    const overview = await loadTemplateOverview();
    console.log(`Loaded the Template overview`);
    setTemplateOverview(overview);
  };

  const onSessionChange = () => {
    (async () => {
      const loadedSessionOverview = await loadSessionOverview();
      console.log(`Loaded the session overview`);
      setSessionOverview(loadedSessionOverview);
    })();
  };

  useEffect(() => {
    (async () => {
      await loadModelOverviewData();
      await loadTemplateOverviewData();
      if (sessionOverview == null) {
        onSessionChange();
      }
    })();
  }, []);

  if (sessionOverview == null || modelOverview == null || templateOverview == null) {
    return <div>Loading...</div>;
  }

  return <MainApp
    modelOverview={modelOverview}
    sessionOverview={sessionOverview}
    templateOverview={templateOverview}
    onSessionChange={onSessionChange}
    />;
}
