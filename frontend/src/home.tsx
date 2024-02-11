import { ModelOverview, ModelSettings, PresetOverview, SessionOverview, TemplateOverview } from "./data";
import { SessionEditor } from "./sessioneditor";
import { SessionOverviewList } from "./sessionoverviewlist";
import { navigate } from "raviger";
import { deleteSession, loadSession, newSession } from "./dataloading";
import { useMemo, useState } from "react";
import classNames from "classnames";
import { PreviewRequestContext } from "./filepreviewrequestcontext";
import { FilePreview } from "./filepreview";


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

  const [previewFileUrl, setPreviewFileUrl] = useState<string | null>(null);
  const [previewFileMimetype, setPreviewFileMimetype] = useState<string | null>(null);

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

  const onPreviewFileUrlRequest = useMemo(() => {
    return (fileUrl: string, fileMimetype:string): void => {
      setPreviewFileUrl(fileUrl);
      setPreviewFileMimetype(fileMimetype);
    };
  }, [sessionId]);

  const closePreview = () => {
    setPreviewFileUrl(null);
    setPreviewFileMimetype(null);
  };

  return (
    <>
    <PreviewRequestContext.Provider value={onPreviewFileUrlRequest}>
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
    </PreviewRequestContext.Provider>

    {previewFileUrl != null &&
      <FilePreview
        fileUrl={previewFileUrl}
        mimetype={previewFileMimetype}
        onCloseRequest={closePreview}
      />
    }
    </>
  );
}
