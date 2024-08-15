import { createViewer3D } from "./viewer"
import Alpine from 'alpinejs'


Alpine.store('lib', () => ({
    models: [],
    viewer: null,
    init() {
        Alpine.store('context').current = 'lib'
        if (!this.$refs.libGlobSidebar) {
            console.log('no sidebar')
            document.getElementsByTagName('body')[0]
                .dispatchEvent(new CustomEvent("load-side", {
                    detail: { context: 'lib' }
                }))
        } else {
            console.log('sidebar exists')
            this.viewer = createViewer3D(this.$refs.viewerContainer);
            this.viewer.setModels(this.models.map(m => { return { id: m } }));
        }
    },
}))