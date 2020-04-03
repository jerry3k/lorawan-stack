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
import { Container, Col, Row } from 'react-grid-system'
import bind from 'autobind-decorator'
import { defineMessages } from 'react-intl'

import PropTypes from '../../../lib/prop-types'
import PageTitle from '../../../components/page-title'
import { withBreadcrumb } from '../../../components/breadcrumbs/context'
import Breadcrumb from '../../../components/breadcrumbs/breadcrumb'
import WebhookForm from '../../components/webhook-form'
import WebhookTemplateForm from '../../components/webhook-template-form'
import sharedMessages from '../../../lib/shared-messages'
import api from '../../api'
import connect from './connect'

const m = defineMessages({
  addCustomWebhook: 'Add custom webhook',
  addWebhookViaTemplate: 'Add webhook from template',
  customWebhook: 'Custom webhook',
})

const breadcrumbs = Component =>
  withBreadcrumb('apps.single.integrations.webhooks.add.form', function(props) {
    const { appId, templateId, webhookTemplate: { name } = {}, isCustom } = props

    return (
      <Breadcrumb
        path={`/applications/${appId}/integrations/webhooks/add/${templateId}`}
        content={isCustom ? m.customWebhook : name}
      />
    )
  })(Component)

@bind
class ApplicationWebhookAddForm extends Component {
  static propTypes = {
    appId: PropTypes.string.isRequired,
    isCustom: PropTypes.bool.isRequired,
    navigateToList: PropTypes.func.isRequired,
    templateId: PropTypes.string,
    webhookTemplate: PropTypes.webhookTemplate,
  }

  static defaultProps = {
    templateId: undefined,
    webhookTemplate: undefined,
  }

  async handleSubmit(webhook) {
    const { appId } = this.props

    await api.application.webhooks.create(appId, webhook)
  }

  handleSubmitSuccess() {
    const { navigateToList, appId } = this.props

    navigateToList(appId)
  }

  render() {
    const { templateId, isCustom, webhookTemplate, appId } = this.props

    let pageTitle

    if (!templateId) {
      pageTitle = sharedMessages.addWebhook
    } else if (isCustom) {
      pageTitle = m.addCustomWebhook
    } else {
      pageTitle = m.addWebhookViaTemplate
    }

    return (
      <Container>
        <PageTitle title={pageTitle} />
        <Row>
          <Col lg={8} md={12}>
            {isCustom ? (
              <WebhookForm
                update={false}
                onSubmit={this.handleSubmit}
                onSubmitSuccess={this.handleSubmitSuccess}
              />
            ) : (
              <WebhookTemplateForm
                appId={appId}
                templateId={templateId}
                onSubmit={this.handleSubmit}
                onSubmitSuccess={this.handleSubmitSuccess}
                webhookTemplate={webhookTemplate}
              />
            )}
          </Col>
        </Row>
      </Container>
    )
  }
}

export default connect(breadcrumbs(ApplicationWebhookAddForm))
export const ApplicationWebhookAddFormWithoutBreadcrumbs = connect(ApplicationWebhookAddForm)
