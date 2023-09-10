import { useEffect, useState } from "react";
import { ModelOverview, PresetOverview, SessionOverview, TemplateOverview } from "./data";
import { loadModelOverview, loadPresetOverview, loadSessionOverview, loadTemplateOverview, scanModels
} from "./dataloading";
import { MainApp } from "./mainapp";

export function LoadingGate() {
  const [sessionOverview, setSessionOverview] = useState<SessionOverview | null>(null);
  const [modelOverview, setModelOverview] = useState<ModelOverview | null>(null);
  const [templateOverview, setTemplateOverview] = useState<TemplateOverview | null>(null);
  const [presetOverview, setPresetOverview] = useState<PresetOverview | null>(null);

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

  const loadPresetOverviewData = async () => {
    const overview = await loadPresetOverview();
    console.log(`Loaded the Preset overview`);
    setPresetOverview(overview);
  };

  const onSessionChange = () => {
    (async () => {
      const loadedSessionOverview = await loadSessionOverview();
      console.log(`Loaded the session overview`);
      setSessionOverview(loadedSessionOverview);
    })();
  };

  const rescanModels = async () => {
    const overview = await scanModels();
    setModelOverview(overview);
  };

  useEffect(() => {
    (async () => {
      await loadModelOverviewData();
      await loadTemplateOverviewData();
      await loadPresetOverviewData();
      if (sessionOverview == null) {
        onSessionChange();
      }
    })();
  }, []);

  if (sessionOverview == null || presetOverview == null || modelOverview == null || templateOverview == null) {
    return <div>Loading...</div>;
  }

  return <MainApp
    modelOverview={modelOverview}
    presetOverview={presetOverview}
    rescanModels={rescanModels}
    sessionOverview={sessionOverview}
    templateOverview={templateOverview}
    onSessionChange={onSessionChange}
    />;
}
