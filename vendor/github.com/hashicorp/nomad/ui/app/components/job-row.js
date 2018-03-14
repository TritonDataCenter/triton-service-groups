import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { lazyClick } from '../helpers/lazy-click';
import { watchRelationship } from 'nomad-ui/utils/properties/watch';
import WithVisibilityDetection from 'nomad-ui/mixins/with-component-visibility-detection';

export default Component.extend(WithVisibilityDetection, {
  store: service(),

  tagName: 'tr',
  classNames: ['job-row', 'is-interactive'],

  job: null,

  onClick() {},

  click(event) {
    lazyClick([this.get('onClick'), event]);
  },

  didReceiveAttrs() {
    // Reload the job in order to get detail information
    const job = this.get('job');
    if (job && !job.get('isLoading')) {
      job.reload().then(() => {
        this.get('watch').perform(job, 100);
      });
    }
  },

  visibilityHandler() {
    if (document.hidden) {
      this.get('watch').cancelAll();
    } else {
      const job = this.get('job');
      if (job && !job.get('isLoading')) {
        this.get('watch').perform(job, 100);
      }
    }
  },

  willDestroy() {
    this.get('watch').cancelAll();
    this._super(...arguments);
  },

  watch: watchRelationship('summary'),
});
