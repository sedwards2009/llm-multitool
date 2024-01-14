import { ModelOverview, ModelSettings, PresetOverview, SessionOverview, TemplateOverview } from "./data";
import { SessionEditor } from "./sessioneditor";
import { SessionOverviewList } from "./sessionoverviewlist";
import { navigate } from "raviger";
import { deleteSession, loadSession, newSession } from "./dataloading";
import { useState } from "react";
import classNames from "classnames";


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

  const [isCreatingSession, setIsCreatingSession] = useState<boolean>(false);

  const onSessionDelete = () => {
    (async () => {
      await deleteSession(sessionId);
      onSessionChange();
      navigate("/");
    })();
  };

  const onNewSessionClick = () => {
    setIsCreatingSession(true);
    (async () => {

      let previousSettings: ModelSettings | null = null;
      if (sessionId != null) {
        const loadedSession = await loadSession(sessionId);
        if (loadedSession != null) {
          previousSettings = loadedSession.modelSettings;
        }
      }

      const session = await newSession(previousSettings);
      setIsCreatingSession(false);
      if (session == null) {
        console.log(`Unable to create a new session.`);
      } else {
        onSessionChange();
        navigate(`/session/${session.id}`);
      }
    })();
  };

  return (
    <div className="top-layout">
      <div className="session-list">
        <button
          className={classNames({"primary": !isCreatingSession})}
          disabled={isCreatingSession}
          onClick={onNewSessionClick}>
            {isCreatingSession ? "Creating session..." : "New Session" }
        </button>
        <p></p>
        <SessionOverviewList
          sessionOverview={sessionOverview}
          selectedSessionId={sessionId}
        />
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
