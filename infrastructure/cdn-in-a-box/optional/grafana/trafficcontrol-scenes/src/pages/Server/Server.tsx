import React, { useMemo } from 'react';
import { SceneApp, SceneAppPage } from '@grafana/scenes';
import { ROUTES } from '../../constants';
import { prefixRoute } from '../../utils/utils.routing';
import { getServerScene } from './scene';

const getScene = () =>
  new SceneApp({
    pages: [
      new SceneAppPage({
        title: 'Server',
        url: prefixRoute(`${ROUTES.Server}`),
        hideFromBreadcrumbs: true,
        getScene: getServerScene,
      }),
    ],
  });

export const ServerPage = () => {
  const scene = useMemo(() => getScene(), []);

  return <scene.Component model={scene} />;
};
