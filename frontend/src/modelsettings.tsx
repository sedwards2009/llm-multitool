import { } from "react";
import { ModelOverview, PresetOverview, TemplateOverview } from "./data";

export interface Props {
  modelOverview: ModelOverview;
  templateOverview: TemplateOverview;
  presetOverview: PresetOverview;
  selectedModelId: string | null;
  setSelectedModelId: (id: string) => void;
  selectedTemplateId: string | null;
  setSelectedTemplateId: (id: string) => void;
  selectedPresetId: string | null;
  setSelectedPresetId: (id: string) => void;
}

export function ModelSettings({modelOverview, presetOverview, templateOverview, selectedModelId, setSelectedModelId,
    selectedTemplateId, setSelectedTemplateId, selectedPresetId, setSelectedPresetId}: Props): JSX.Element  {

  const modelIDs: (string | null)[] = modelOverview.models.map(m => m.id);
  const templateIDs: (string | null)[] = templateOverview.templates.map(t => t.id);
  const presetIDs: (string | null)[] = presetOverview.presets.map(p => p.id);

  return <div className="gui-layout cols-1-2">
    <div>
      <i className="fa fa-robot"></i>&nbsp;&nbsp;Model:
    </div>
    <div>
      <select
        value={selectedModelId == null ? undefined : selectedModelId}
        onChange={(e) => setSelectedModelId(e.target.value)}>
        {
          !modelIDs.includes(selectedModelId) && <option key={""} value="">(none)</option>
        }
        {
          modelOverview.models.map(m => <option key={m.id} value={m.id}>{ m.name }</option>)
        }
      </select>
    </div>
    <div>
      <i className="fa fa-hammer"></i>&nbsp;&nbsp;Task:
    </div>
    <div>
      <select
        value={selectedTemplateId == null ? undefined : selectedTemplateId}
        onChange={(e) => setSelectedTemplateId(e.target.value)}>
        {
          !templateIDs.includes(selectedTemplateId) && <option key={""} value="">(none)</option>
        }
        {
          templateOverview.templates.map(t => <option key={t.id} value={t.id}>{ t.name }</option>)
        }
      </select>
    </div>
    <div>
      <i className="fa fa-tachometer-alt"></i>&nbsp;&nbsp;Creativeness:
    </div>
    <div>
      <select
        value={selectedPresetId == null ? undefined : selectedPresetId}
        onChange={(e) => setSelectedPresetId(e.target.value)}>
        {
          !presetIDs.includes(selectedPresetId) && <option key={""} value="">(none)</option>
        }
        {
          presetOverview.presets.map(p => <option key={p.id} value={p.id}>{ p.name }</option>)
        }
      </select>
    </div>
  </div>;
}