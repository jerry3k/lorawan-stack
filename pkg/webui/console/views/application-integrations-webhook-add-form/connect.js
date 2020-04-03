// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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

import { connect } from 'react-redux'
import { push } from 'connected-react-router'

import { getWebhookTemplate } from '../../store/actions/webhook-templates'
import { selectSelectedApplicationId } from '../../store/selectors/applications'
import { selectWebhookTemplateById } from '../../store/selectors/webhook-templates'

export default ApplicationWebhookAddForm =>
  connect(
    (state, props) => {
      const templateId = props.match.params.templateId
      return {
        appId: selectSelectedApplicationId(state),
        templateId,
        isCustom: !templateId || templateId === 'custom',
        webhookTemplate: selectWebhookTemplateById(state, templateId),
      }
    },
    dispatch => ({
      getWebhookTemplate: (templateId, selector) =>
        dispatch(getWebhookTemplate(templateId, selector)),
      navigateToList: appId => dispatch(push(`/applications/${appId}/integrations/webhooks`)),
    }),
  )(ApplicationWebhookAddForm)
