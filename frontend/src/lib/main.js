import { createViewer3D } from "./viewer"
import Alpine from 'alpinejs'


Alpine.data('lib', () => ({
    tab: '',
    models: [],
    viewer: null,
    init() {
        Alpine.store('context').current = 'lib'
        if (this.$refs.libSidebar) {
            console.log('sidebar exists')
            this.viewer = createViewer3D(this.$refs.viewerContainer);
            this.viewer.setModels(this.models.map(m => { return { id: m } }));
        }
    },
    openTab(tab, context) {
        if (this.tab === tab && !context) {
            this.tab = ''
            return
        }
        var assetID = null
        if (context) {
            assetID = context.assetID
        } else {
            assetID = document.getElementsByName("assetID")[0].value
        }

        document.getElementsByTagName("body")[0].dispatchEvent(new CustomEvent('tab-' + tab, { detail: { assetID } }))
        this.tab = tab
    },
    addModel(model) {
        this.models.push(model)
        this.viewer.setModels(this.models.map(m => { return { id: m } }));
        this.tab = 'viewer'
    },
}))