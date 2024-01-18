import React, { useMemo } from 'react';
import { SceneApp, SceneAppPage } from '@grafana/scenes';
import { ROUTES } from '../../constants';
import { prefixRoute } from '../../utils/utils.routing';
import { getCacheGroupScene } from './scene';

const getScene = () =>
  new SceneApp({
    pages: [
      new SceneAppPage({
        title: 'Cache Groups',
        url: prefixRoute(`${ROUTES.CacheGroup}`),
        hideFromBreadcrumbs: true,
        getScene: getCacheGroupScene,
      }),
    ],
  });

export const CacheGroupPage = () => {
  const scene = useMemo(() => getScene(), []);

  return <scene.Component model={scene} />;
};
