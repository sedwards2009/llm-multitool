import { useRoutes } from "raviger";
import { ModelOverview, PresetOverview, SessionOverview, TemplateOverview } from "./data";
import { EmptyHome } from "./emptyhome";
import { Home } from "./home";
import { SettingsPage } from "./settingspage";
import { TitleBar } from "./titlebar";

export interface Props {
  modelOverview: ModelOverview;
  rescanModels: () => Promise<void>;
  presetOverview: PresetOverview;
  sessionOverview: SessionOverview;
  templateOverview: TemplateOverview;
  onSessionChange: () => void;
}

export function MainApp({ modelOverview, presetOverview, sessionOverview, templateOverview, onSessionChange,
  rescanModels }: Props): JSX.Element {

  return (
    <>
      {
        useRoutes(
          {
            '/': () => {
              return <>
                <TitleBar />
                <EmptyHome sessionOverview={sessionOverview} onSessionChange={onSessionChange}/>
              </>;
            },
            '/session/:sessionId': ({ sessionId }: { sessionId: any }) => {
              return <>
                <TitleBar />
                <Home
                  modelOverview={modelOverview}
                  presetOverview={presetOverview}
                  sessionOverview={sessionOverview}
                  templateOverview={templateOverview}
                  sessionId={sessionId}
                  onSessionChange={onSessionChange}
                />
              </>;
            },
            '/settings': () => {
              return <SettingsPage
                modelOverview={modelOverview}
                rescanModels={rescanModels}
              />
            }
          }
        )
      }
    </>
  );
}
