import { ModelOverview, PresetOverview, SessionOverview, TemplateOverview } from "./data";
import { SessionEditor } from "./sessioneditor";
import { SessionOverviewList } from "./sessionoverviewlist";
import { NewSessionButton } from "./newsessionbutton";
import { navigate } from "raviger";
import { deleteSession } from "./dataloading";

export interface Props {
  modelOverview: ModelOverview;
  presetOverview: PresetOverview;
  sessionOverview: SessionOverview;
  templateOverview: TemplateOverview;
  sessionId: string;
  onSessionChange: ()=> void;
}

export function Home({ modelOverview, presetOverview, sessionOverview, templateOverview, sessionId, onSessionChange
    }: Props): JSX.Element {

  const onSessionDelete = () => {
    (async () => {
      await deleteSession(sessionId);
      onSessionChange();
      navigate("/");
    })();
  };

  return (
    <div className="top-layout">
      <div className="session-list">
        <SessionOverviewList
          sessionOverview={sessionOverview}
          selectedSessionId={sessionId}
        />
        <NewSessionButton onSessionChange={onSessionChange} />
      </div>
      <div className="session-tab">
        <SessionEditor
          key={sessionId}
          sessionId={sessionId}
          modelOverview={modelOverview}
          presetOverview={presetOverview}
          templateOverview={templateOverview}
          onSessionDelete={onSessionDelete}
          onSessionChange={onSessionChange}
        />
      </div>
  </div>
  );
}
