import Ember from 'ember';
import notifyError from 'nomad-ui/utils/notify-error';

const { Route, RSVP, inject } = Ember;

export default Route.extend({
  store: inject.service(),

  serialize(model) {
    return { job_name: model.get('plainId') };
  },

  model(params, transition) {
    const namespace = transition.queryParams.namespace || this.get('system.activeNamespace.id');
    const name = params.job_name;
    const fullId = JSON.stringify([name, namespace]);
    return this.get('store')
      .findRecord('job', fullId, { reload: true })
      .then(job => {
        return RSVP.all([job.get('allocations'), job.get('evaluations')]).then(() => job);
      })
      .catch(notifyError(this));
  },
});
