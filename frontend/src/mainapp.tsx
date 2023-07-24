import { useRoutes } from "raviger";
import { ModelOverview, SessionOverview } from "./data";
import { EmptyHome } from "./emptyhome";
import { Home } from "./home";

export interface Props {
  modelOverview: ModelOverview;
  sessionOverview: SessionOverview;
  onSessionChange: () => void;
}

export function MainApp({ modelOverview, sessionOverview, onSessionChange }: Props): JSX.Element {
  return (
    <>
      <h1>LLM Workbench</h1>
      {
        useRoutes(
          {
            '/': () => <EmptyHome sessionOverview={sessionOverview} onSessionChange={onSessionChange}/>,
            '/session/:sessionId': ({ sessionId }: { sessionId: any }) => {
              return <Home
                modelOverview={modelOverview}
                sessionOverview={sessionOverview}
                sessionId={sessionId}
                onSessionChange={onSessionChange}
              />;
            }
          }
        )
      }
    </>
  );
}
