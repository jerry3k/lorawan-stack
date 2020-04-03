// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React, { Component } from 'react'
import bind from 'autobind-decorator'
import { connect } from 'react-redux'
import { Switch, Route } from 'react-router'

import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import ApplicationWebhookAddChoose from '../application-integrations-webhook-add-choose'
import ApplicationWebhookAddForm, {
  ApplicationWebhookAddFormWithoutBreadcrumbs,
} from '../application-integrations-webhook-add-form'
import { selectSelectedApplicationId } from '../../store/selectors/applications'
import { selectWebhookTemplates } from '../../store/selectors/webhook-templates'
import sharedMessages from '../../../lib/shared-messages'
import PropTypes from '../../../lib/prop-types'

@connect(state => ({
  appId: selectSelectedApplicationId(state),
  hasTemplates: selectWebhookTemplates(state).length !== 0,
}))
@withBreadcrumb('apps.single.integrations.webhooks.add', function(props) {
  const { appId } = props
  return (
    <Breadcrumb
      path={`/applications/${appId}/integrations/webhooks/add`}
      content={sharedMessages.add}
    />
  )
})
@bind
export default class ApplicationWebhookAdd extends Component {
  static propTypes = {
    hasTemplates: PropTypes.bool.isRequired,
    match: PropTypes.match.isRequired,
  }
  render() {
    const {
      match,
      match: { url: path },
      hasTemplates,
    } = this.props

    // Do not render the chooser when there are no webhook templates configured
    if (!hasTemplates) {
      return <ApplicationWebhookAddFormWithoutBreadcrumbs match={match} />
    }

    return (
      <Switch>
        <Route exact path={`${path}`} component={ApplicationWebhookAddChoose} />
        <Route path={`${path}/:templateId`} component={ApplicationWebhookAddForm} />
      </Switch>
    )
  }
}
