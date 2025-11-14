import axios from 'axios';

import { BACKEND_HTTP_URL_ROOT } from 'utils/constants';

// const getSessions = async () => {
//   const { data } = await axios.get('/settings.json');
//   console.log(data);
// };

// getSessions();
const backendAxiosInstance = axios.create({
  baseURL: BACKEND_HTTP_URL_ROOT,
});

export default backendAxiosInstance;
