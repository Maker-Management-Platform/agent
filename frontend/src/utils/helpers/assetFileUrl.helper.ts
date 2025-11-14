import { BACKEND_HTTP_URL_ROOT } from 'utils/constants';

const assetFileUrl = (uuid: string, default_image_id: string) => (
  [
    BACKEND_HTTP_URL_ROOT,
    'projects',
    uuid,
    'assets',
    default_image_id,
    'file',
  ].join('/')
);

export default assetFileUrl;
