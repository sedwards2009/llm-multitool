import { useState } from "react";
import { navigate } from "raviger";
import classNames from "classnames";
import { newSession } from "./dataloading";

export interface Props {
  onSessionChange: ()=> void;
}

export function NewSessionButton({ onSessionChange }: Props): JSX.Element {
  const [isCreatingSession, setIsCreatingSession] = useState<boolean>(false);

  const onClick = () => {
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

  return <button
    className={classNames({"primary": !isCreatingSession})}
    disabled={isCreatingSession}
    onClick={onClick}>
      {isCreatingSession ? "Creating session..." : "New Session" }
  </button>;
}
