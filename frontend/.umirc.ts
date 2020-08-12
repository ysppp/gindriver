import { IConfig } from 'umi-types'; // ref: https://umijs.org/config/

const config: IConfig = {
  routes: [
    {
      path: '/',
      component: '../layouts/index',
      routes: [
        {
          path: '/user',
          component: './user',
        },
        {
          path: '/',
          component: '../pages/index',
        },
      ],
    },
  ],
  // plugins: [
  //   // ref: https://umijs.org/plugin/umi-plugin-react.html
  //   [
  //     'umi-plugin-react',
  //     {
  //       antd: true,
  //       dva: false,
  //       dynamicImport: false,
  //       title: 'frontend',
  //       dll: false,
  //       routes: {
  //         exclude: [/components\//],
  //       },
  //     },
  //   ],
  // ],
  antd: {},
  title: "GinDriver"
};

export default config;
