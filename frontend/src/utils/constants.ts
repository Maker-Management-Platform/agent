const BACKEND_HTTP_URL_ROOT = [
  import.meta.env['MMP_BACKEND_HOST'] ?? globalThis.location.host,
  import.meta.env['MMP_BACKEND_PATH'],
].join('');

const PROJECTS_PER_PAGE = [
  10,
  20,
  50,
  100,
];

export {
  BACKEND_HTTP_URL_ROOT,
  PROJECTS_PER_PAGE,
};
