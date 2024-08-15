import { createViewer3D } from "./viewer"

var viewer = null;
var models = [];
var viewerEl = null;
(function () {
    viewerEl = document.getElementById('viewer')
    if (viewerEl) {
        setControlEvents(viewerEl);
        viewer = createViewer3D(viewerEl.querySelector('.viewer-container'));
        viewer.setModels(models.map(m => { return { id: m } }));
    }
    document.body.addEventListener('htmx:load', function (evt) {
        var newViewerEl = evt.detail.elt.querySelector('.viewer-container');
        if (newViewerEl) {
            console.log('newViewerEl', newViewerEl)
            if (viewer) {
                console.log('destroying viewer')
                //viewer.destroy();
            }
            // initialize buttons for the view wrapper
            console.log('creating viewer')
            //viewer = createViewer3D(newViewerEl);
            //viewer.setModels(models.map(m => { return { id: m } }));
        }
        addListeners();
    });
})()

function addListeners() {
    document.querySelectorAll('button.add-model').forEach((el) => {
        el.addEventListener('click', function (e) {
            var m = e.currentTarget.dataset.modelId;
            if (models.indexOf(m) === -1) {
                models.push(m);
            }
            viewerEl.classList.remove('hidden');
            viewer.setModels(models.map(m => { return { id: m } }));
        })
    })
}

function setControlEvents(viewerEl) {
    if (!viewerEl) {
        return;
    }
    viewerEl.querySelector('.close').addEventListener('click', function (e) {
        viewerEl.classList.add('hidden');
    })
}
