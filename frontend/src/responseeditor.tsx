import { Response } from "./data";

export interface Props {
  response: Response;
}

export function ResponseEditor({response}: Props): JSX.Element {
  return <div className="card char-width-20">
    <h4>Prompt:</h4>
    {response.prompt}<br />
    <h4>Output:</h4>
    {response.text}
  </div>;
}
