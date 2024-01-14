import { useEffect, useState } from "react";
import classNames from "classnames";

import { SessionOverview } from "./data";
import { navigate } from "raviger";
import { newSession } from "./dataloading";

export interface Props {
  sessionOverview: SessionOverview;
  onSessionChange: ()=> void;
}

export function EmptyHome({ sessionOverview, onSessionChange }: Props): JSX.Element {
  useEffect(() => {
    if (sessionOverview.sessionSummaries.length === 0) {
      return;
    }
    navigate(`/session/${sessionOverview.sessionSummaries[0].id}`);
  }, [sessionOverview]);

  const [isCreatingSession, setIsCreatingSession] = useState<boolean>(false);

  const onNewSessionClick = () => {
    setIsCreatingSession(true);
    (async () => {
      const session = await newSession();
      setIsCreatingSession(false);
      if (session == null) {
        console.log(`Unable to create a new session.`);
      } else {
        onSessionChange();
        navigate(`/session/${session.id}`);
      }
    })();
  };

  return <div>
    <button
      className={classNames({"primary": !isCreatingSession})}
      disabled={isCreatingSession}
      onClick={onNewSessionClick}>
        {isCreatingSession ? "Creating session..." : "New Session" }
    </button>
  </div>;
}
