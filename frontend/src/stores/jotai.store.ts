import { createStore } from 'jotai';

// Create store outside of provider to allow access outside of React components.
const jotaiStore = createStore();

export default jotaiStore;
