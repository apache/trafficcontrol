import React from 'react';

import { SceneComponentProps, SceneObjectBase, SceneObjectState } from '@grafana/scenes';
import { InlineField, Input } from '@grafana/ui';

interface ServerState extends SceneObjectState {
  name: string;
}

function Renderer({ model }: SceneComponentProps<ServerCustomObject>) {
  const { name } = model.useState();

  return (
    <InlineField label="server" style={{ margin: '0' }}>
      <Input
        prefix=""
        defaultValue={name}
        width={20}
        type="string"
        onBlur={(evt) => {
          model.onValueChange(evt.currentTarget.value);
        }}
      />
    </InlineField>
  );
}

export class ServerCustomObject extends SceneObjectBase<ServerState> {
  public static Component = Renderer;

  onValueChange = (value: string) => {
    this.setState({ name: value });
  };
}
