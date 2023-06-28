import { Session } from "./data";

export interface Props {
  session: Session;
}

export function SessionEditor({session}: Props): JSX.Element {
  return <div className="session-editor">
    <div className="session-prompt-pane">
      <h3>Prompt</h3>
      <textarea defaultValue={session.prompt} /><br />
      <button className="success">Submit</button>
    </div>
    <div className="session-response-pane">
      <h3>Responses</h3>
    </div>
  </div>;
}
