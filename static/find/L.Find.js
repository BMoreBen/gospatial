/**
 * FIND Client
 * Author: Stefan Safranek
 * Email:  sjsafranek@gmail.com
 */

L.Find = L.Class.extend({

    options: {},

    initialize: function(datasources, options) {
        L.setOptions(this, options || {});
        this._map = null;
        this.datasources = datasources;
        this.featureLayers = {};
    },

    addTo: function(map) {
        find = this;
        this._map = map;
        this.addUiControls();
        this.getLayer($('#layers').val());
        this._map.fitBounds(
            find.featureLayers[$('#layers').val()].getBounds()
        );
        return this;
    },

    preventPropogation: function(obj) {
        map = this._map;
        // http://gis.stackexchange.com/questions/104507/disable-panning-dragging-on-leaflet-map-for-div-within-map
        // Disable dragging when user's cursor enters the element
        obj.getContainer().addEventListener('mouseover', function () {
            map.dragging.disable();
        });
        // Re-enable dragging when user's cursor leaves the element
        obj.getContainer().addEventListener('mouseout', function () {
            map.dragging.enable();
        });
    },

    addUiControls: function() {
        this._addLogoControl();
        this._addLayerControl();
        this._addMouseControl();
        this._addFeatureAttributesControl(); 
    },

    _addLogoControl: function() {
        var logo = L.control({position : 'topleft'});
        logo.onAdd = function () {
            // this._div = L.DomUtil.create('div', 'logo-compass');
            // this._div.innerHTML = "<img class='img-logo-compass' src='/images/compass.png' alt='logo'>"
            this._div = L.DomUtil.create('div', 'logo-hypercube');
            this._div.innerHTML = "<img class='img-logo-hypercube' src='/images/HyperCube2.png' alt='logo'>"
            return this._div;
        };
        logo.addTo(this._map);
    },

    _addLayerControl: function() {
        find = this;
        // Create UI control element
        geojsonLayerControl = L.control({position: 'topright'});
        geojsonLayerControl.onAdd = function () {
            var div = L.DomUtil.create('div', 'info legend');
            div.innerHTML = '';
            div.innerHTML += '<i class="fa fa-search-plus" id="zoom" style="padding-left:5px; margin-right:0px;"></i><select name="geojson" id="layers"></select>';
            return div;
        };
        geojsonLayerControl.addTo(this._map);
        this.preventPropogation(geojsonLayerControl);
        // Fill drop down options
        for (var _i=0; _i < this.datasources.length; _i++) {
            var obj = document.createElement('option');
            obj.value = this.datasources[_i];
            obj.text = this.datasources[_i];
            $('#layers').append(obj);
        }
        // UI Events listeners
        $('#layers').on('change', function(){ 
            find.getLayer($('#layers').val());
        });
        $('#zoom').on('click', function(){ 
            find._map.fitBounds(
                find.featureLayers[$('#layers').val()].getBounds()
            );
        });
    },

    _addMouseControl: function() {
        // Create UI control element
        mouseLocationControl = L.control({position: 'bottomright'});
        mouseLocationControl.onAdd = function () {
            var div = L.DomUtil.create('div');
            div.innerHTML = "<div id='location'></div>";
            return div;
        };
        mouseLocationControl.addTo(this._map);
        this.preventPropogation(mouseLocationControl);
        // UI Event listeners
        this._map.on('mousemove', function(e) {
            $("#location")[0].innerHTML = "<strong>Lat, Lon : " + e.latlng.lat.toFixed(4) + ", " + e.latlng.lng.toFixed(4) + "</strong>";
        });
    },

    _addFeatureAttributesControl: function() {
        // Create UI control element
        featureAttributesControl = L.control({position: 'bottomright'});
        featureAttributesControl.onAdd = function () {
            var div = L.DomUtil.create('div', 'info legend');
            div.innerHTML = "<div id='attributes'>Hover over features</div>";
            return div;
        };
        featureAttributesControl.addTo(this._map);
        this.preventPropogation(featureAttributesControl);
    },

    getLayer: function(datasource) {
        data = this.getRequest("/api/v1/layer/" + datasource);
        this.updateFeatureLayers(data);
    },

    updateFeatureLayers: function(data) {
        for (var _i in this.featureLayers){
            if (this._map.hasLayer(this.featureLayers[_i])) {
                this._map.removeLayer(this.featureLayers[_i]);
            }
        }
        try {
            this.featureLayers[$('#layers').val()] = this.createFeatureLayer(data);
            this.featureLayers[$('#layers').val()].addTo(this._map);
        }
        catch(err) { console.log(err); }
    },

    createFeatureLayer: function(data) {
        var featureLayer = L.geoJson(data, {
            style: {
                "weight": 2, 
                "color": "#000", 
                "fillOpacity": 0.25,
            },
            pointToLayer: function(feature, latlng) {
                return L.circleMarker(latlng, {
                    radius: 4,
                    weight: 1,
                    fillOpacity: 0.25,
                    color: '#000'
                });
            },
            onEachFeature: function (feature, layer) {

                function highlightFeature(e) {
                    var layer = e.target;
                    layer.setStyle({
                        weight: 5,
                        color: '#000',
                        dashArray: '',
                        fillOpacity: 0.5
                    });
                    if (!L.Browser.ie && !L.Browser.opera) {
                        layer.bringToFront();
                    }
                }

                function resetHighlight(e) {
                    featureLayer.resetStyle(e.target);
                }

                function zoomToFeature(e) {
                    map.fitBounds(e.target.getBounds());
                }
                
                // layer.bindPopup(
                //     "<button onclick=map.editfeature(" + JSON.stringify(feature) + ")>Edit</button>"
                // );
                layer.on({
                    mouseover: function(feature){
                        var properties = feature.target.feature.properties;
                        var results = "<table>";
                        results += "<th>Field</th><th>Attribute</th>";
                        for (var item in properties) {
                            results += "<tr><td>" + item + "</td><td>" + properties[item] + "</td></tr>";
                        }
                        results += "</table>";
                        $("#attributes")[0].innerHTML = results;
                        highlightFeature(feature);
                    },
                    mouseout: function(e){
                        $("#attributes")[0].innerHTML = "Hover over features";
                        resetHighlight(e);
                    },
                    click: function(feature) {
                        var properties = feature.target.feature.properties;
                        var results = "";
                        for (var item in properties) {
                            results += item + ": " + properties[item] + "<br>";
                        }
                        layer.bindPopup(results);
                    }
                });
            }
        });
        return featureLayer;
    },

    getFeature: function(datasource, k){
        results = this.getRequest("/api/v1/layer/" + datasource + "/feature/" + k);
        return results;
    },

    postRequest: function(route, data) {
        var results;
        $.ajax({
            crossDomain: true,
            type: "POST",
            async: false,
            data: data,
            url: route,
            dataType: 'JSON',
            success: function (data) {
                try {
                    results = data;
                }
                catch(err){  console.log('Error:', err);  }
            },
            error: function(xhr,errmsg,err) {
                console.log(xhr.status,xhr.responseText,errmsg,err);
                result = null;
            }
        });
        return results;
    },

    getRequest: function(route, data) {
        var results;
        $.ajax({
            crossDomain: true,
            type: "GET",
            async: false,
            data: data,
            url: route,
            dataType: 'JSON',
            success: function (data) {
                try {
                    results = data;
                }
                catch(err){  console.log('Error:', err);  }
            },rror: function(xhr,errmsg,err) {
                console.log(xhr.status,xhr.responseText,errmsg,err);
                result = null;
            }
        });
        return results;
    }

});

L.find = function(datasources, options) {
    return new L.Find(datasources, options);
};