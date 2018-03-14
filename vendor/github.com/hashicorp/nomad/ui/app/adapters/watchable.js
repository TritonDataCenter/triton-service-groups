import { get, computed } from '@ember/object';
import { assign } from '@ember/polyfills';
import { inject as service } from '@ember/service';
import queryString from 'npm:query-string';
import ApplicationAdapter from './application';
import { AbortError } from 'ember-data/adapters/errors';

export default ApplicationAdapter.extend({
  watchList: service(),
  store: service(),

  xhrs: computed(function() {
    return {};
  }),

  ajaxOptions(url) {
    const ajaxOptions = this._super(...arguments);

    const previousBeforeSend = ajaxOptions.beforeSend;
    ajaxOptions.beforeSend = function(jqXHR) {
      if (previousBeforeSend) {
        previousBeforeSend(...arguments);
      }
      this.get('xhrs')[url] = jqXHR;
      jqXHR.always(() => {
        delete this.get('xhrs')[url];
      });
    };

    return ajaxOptions;
  },

  findAll(store, type, sinceToken, snapshotRecordArray, additionalParams = {}) {
    const params = assign(this.buildQuery(), additionalParams);
    const url = this.urlForFindAll(type.modelName);

    if (get(snapshotRecordArray || {}, 'adapterOptions.watch')) {
      params.index = this.get('watchList').getIndexFor(url);
      this.cancelFindAll(type.modelName);
    }

    return this.ajax(url, 'GET', {
      data: params,
    });
  },

  findRecord(store, type, id, snapshot, additionalParams = {}) {
    let [url, params] = this.buildURL(type.modelName, id, snapshot, 'findRecord').split('?');
    params = assign(queryString.parse(params) || {}, this.buildQuery(), additionalParams);

    if (get(snapshot || {}, 'adapterOptions.watch')) {
      params.index = this.get('watchList').getIndexFor(url);
      this.cancelFindRecord(type.modelName, id);
    }

    return this.ajax(url, 'GET', {
      data: params,
    }).catch(error => {
      if (error instanceof AbortError) {
        return {};
      }
      throw error;
    });
  },

  reloadRelationship(model, relationshipName, watch = false) {
    const relationship = model.relationshipFor(relationshipName);
    if (relationship.kind !== 'belongsTo' && relationship.kind !== 'hasMany') {
      throw new Error(
        `${relationship.key} must be a belongsTo or hasMany, instead it was ${relationship.kind}`
      );
    } else {
      const url = model[relationship.kind](relationship.key).link();
      let params = {};

      if (watch) {
        params.index = this.get('watchList').getIndexFor(url);
        this.cancelReloadRelationship(model, relationshipName);
      }

      if (url.includes('?')) {
        params = assign(queryString.parse(url.split('?')[1]), params);
      }

      return this.ajax(url, 'GET', {
        data: params,
      }).then(
        json => {
          const store = this.get('store');
          const normalizeMethod =
            relationship.kind === 'belongsTo'
              ? 'normalizeFindBelongsToResponse'
              : 'normalizeFindHasManyResponse';
          const serializer = store.serializerFor(relationship.type);
          const modelClass = store.modelFor(relationship.type);
          const normalizedData = serializer[normalizeMethod](store, modelClass, json);
          store.push(normalizedData);
        },
        error => {
          if (error instanceof AbortError) {
            return relationship.kind === 'belongsTo' ? {} : [];
          }
          throw error;
        }
      );
    }
  },

  handleResponse(status, headers, payload, requestData) {
    const newIndex = headers['x-nomad-index'];
    if (newIndex) {
      this.get('watchList').setIndexFor(requestData.url, newIndex);
    }

    return this._super(...arguments);
  },

  cancelFindRecord(modelName, id) {
    const url = this.urlForFindRecord(id, modelName);
    const xhr = this.get('xhrs')[url];
    if (xhr) {
      xhr.abort();
    }
  },

  cancelFindAll(modelName) {
    const xhr = this.get('xhrs')[this.urlForFindAll(modelName)];
    if (xhr) {
      xhr.abort();
    }
  },

  cancelReloadRelationship(model, relationshipName) {
    const relationship = model.relationshipFor(relationshipName);
    if (relationship.kind !== 'belongsTo' && relationship.kind !== 'hasMany') {
      throw new Error(
        `${relationship.key} must be a belongsTo or hasMany, instead it was ${relationship.kind}`
      );
    } else {
      const url = model[relationship.kind](relationship.key).link();
      const xhr = this.get('xhrs')[url];
      if (xhr) {
        xhr.abort();
      }
    }
  },
});
