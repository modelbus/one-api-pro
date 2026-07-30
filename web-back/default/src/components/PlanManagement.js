import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Button,
  Form,
  Label,
  Modal,
  Table,
} from 'semantic-ui-react';
import {
  API,
  showError,
  showSuccess,
} from '../helpers';

const PlanManagement = () => {
  const { t } = useTranslation();
  const [plans, setPlans] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showEditModal, setShowEditModal] = useState(false);
  const [editForm, setEditForm] = useState({
    id: 0,
    name: '',
    description: '',
    price: 0,
    duration_days: 30,
    duration_text: '',
    status: 1,
    recommended: false,
    sort: 0,
    features: '',
    model_limits: '{}',
    default_model: '',
  });

  const loadPlans = async () => {
    setLoading(true);
    const res = await API.get('/api/plan/?p=0');
    const { success, message, data } = res.data;
    if (success) {
      setPlans(data);
    } else {
      showError(message);
    }
    setLoading(false);
  };

  useEffect(() => {
    loadPlans();
  }, []);

  const savePlan = async () => {
    const isEdit = editForm.id > 0;
    const res = isEdit
      ? await API.put('/api/plan/', editForm)
      : await API.post('/api/plan/', editForm);
    const { success, message } = res.data;
    if (success) {
      showSuccess(isEdit ? t('plan.messages.update_success') : t('plan.messages.add_success'));
      setShowEditModal(false);
      loadPlans();
    } else {
      showError(message);
    }
  };

  const deletePlan = async (id) => {
    const res = await API.delete(`/api/plan/${id}/`);
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('plan.messages.delete_success'));
      loadPlans();
    } else {
      showError(message);
    }
  };

  const toggleStatus = async (plan) => {
    const res = await API.put('/api/plan/', {
      ...plan,
      status: plan.status === 1 ? 0 : 1,
    });
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('plan.messages.update_success'));
      loadPlans();
    } else {
      showError(message);
    }
  };

  const openEditModal = (plan) => {
    if (plan) {
      setEditForm({ ...plan, model_limits: plan.model_limits || '{}', default_model: plan.default_model || '' });
    } else {
      setEditForm({
        id: 0,
        name: '',
        description: '',
        price: 0,
        duration_days: 30,
        duration_text: '',
        status: 1,
        recommended: false,
        sort: 0,
        features: '',
        model_limits: '{}',
        default_model: '',
      });
    }
    setShowEditModal(true);
  };

  return (
    <>
      <Button size='small' onClick={() => openEditModal(null)}>
        {t('plan.buttons.add')}
      </Button>
      <Button size='small' onClick={loadPlans} loading={loading}>
        {t('plan.buttons.refresh')}
      </Button>

      <Table basic='very' compact size='small' style={{ marginTop: 10 }}>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>ID</Table.HeaderCell>
            <Table.HeaderCell>{t('plan.table.name')}</Table.HeaderCell>
            <Table.HeaderCell>{t('plan.table.price')}</Table.HeaderCell>
            <Table.HeaderCell>{t('plan.table.duration_days')}</Table.HeaderCell>
            <Table.HeaderCell>{t('plan.table.default_model')}</Table.HeaderCell>
            <Table.HeaderCell>{t('plan.table.recommended')}</Table.HeaderCell>
            <Table.HeaderCell>{t('plan.table.status')}</Table.HeaderCell>
            <Table.HeaderCell>{t('plan.table.sort')}</Table.HeaderCell>
            <Table.HeaderCell>{t('plan.table.actions')}</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {plans.map((plan) => (
            <Table.Row key={plan.id}>
              <Table.Cell>{plan.id}</Table.Cell>
              <Table.Cell>{plan.name}</Table.Cell>
              <Table.Cell>{plan.price}</Table.Cell>
              <Table.Cell>{plan.duration_days}{plan.duration_text || t('plan.table.days')}</Table.Cell>
              <Table.Cell>{plan.default_model || '-'}</Table.Cell>
              <Table.Cell>{plan.recommended ? '★' : ''}</Table.Cell>
              <Table.Cell>
                {plan.status === 1
                  ? <Label basic color='green'>{t('plan.status.enabled')}</Label>
                  : <Label basic color='grey'>{t('plan.status.disabled')}</Label>}
              </Table.Cell>
              <Table.Cell>{plan.sort}</Table.Cell>
              <Table.Cell>
                <Button size='tiny' onClick={() => openEditModal(plan)}>
                  {t('plan.buttons.edit')}
                </Button>
                <Button size='tiny' onClick={() => toggleStatus(plan)}>
                  {plan.status === 1 ? t('plan.buttons.disable') : t('plan.buttons.enable')}
                </Button>
                <Modal
                  trigger={<Button size='tiny' negative>{t('plan.buttons.delete')}</Button>}
                  header={t('plan.modals.delete_confirm')}
                  content={t('plan.modals.delete_content', { name: plan.name })}
                  actions={[
                    { key: 'cancel', content: t('plan.buttons.cancel'), positive: false },
                    { key: 'delete', content: t('plan.buttons.confirm_delete'), positive: true, onClick: () => deletePlan(plan.id) },
                  ]}
                />
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table>

      <Modal open={showEditModal} onClose={() => setShowEditModal(false)} size='large'>
        <Modal.Header>{editForm.id ? t('plan.modals.edit_title') : t('plan.modals.add_title')}</Modal.Header>
        <Modal.Content>
          <Form>
            <Form.Input
              label={t('plan.form.name')}
              value={editForm.name}
              onChange={(e, { value }) => setEditForm({ ...editForm, name: value })}
            />
            <Form.TextArea
              label={t('plan.form.description')}
              value={editForm.description}
              onChange={(e, { value }) => setEditForm({ ...editForm, description: value })}
            />
            <Form.Group widths='three'>
              <Form.Input
                label={t('plan.form.price')}
                type='number'
                value={editForm.price}
                onChange={(e, { value }) => setEditForm({ ...editForm, price: parseFloat(value) || 0 })}
              />
              <Form.Input
                label={t('plan.form.duration_days')}
                type='number'
                value={editForm.duration_days}
                onChange={(e, { value }) => setEditForm({ ...editForm, duration_days: parseInt(value) || 30 })}
              />
              <Form.Input
                label={t('plan.form.duration_text')}
                value={editForm.duration_text}
                onChange={(e, { value }) => setEditForm({ ...editForm, duration_text: value })}
                placeholder={t('plan.form.duration_text_placeholder')}
              />
            </Form.Group>
            <Form.Group widths='three'>
              <Form.Input
                label={t('plan.form.sort')}
                type='number'
                value={editForm.sort}
                onChange={(e, { value }) => setEditForm({ ...editForm, sort: parseInt(value) || 0 })}
              />
              <Form.Checkbox
                label={t('plan.form.recommended')}
                checked={editForm.recommended}
                onChange={(e, { checked }) => setEditForm({ ...editForm, recommended: checked })}
              />
            </Form.Group>
            <Form.TextArea
              label={t('plan.form.features')}
              value={editForm.features}
              onChange={(e, { value }) => setEditForm({ ...editForm, features: value })}
              placeholder={t('plan.form.features_placeholder')}
            />
            <Form.Input
              label={t('plan.form.default_model')}
              value={editForm.default_model}
              onChange={(e, { value }) => setEditForm({ ...editForm, default_model: value })}
              placeholder={t('plan.form.default_model_placeholder')}
            />
            <Form.TextArea
              label={t('plan.form.model_limits')}
              value={editForm.model_limits}
              onChange={(e, { value }) => setEditForm({ ...editForm, model_limits: value })}
              placeholder={t('plan.form.model_limits_placeholder')}
              rows={10}
            />
          </Form>
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={() => setShowEditModal(false)}>{t('plan.buttons.cancel')}</Button>
          <Button positive onClick={savePlan}>{t('plan.buttons.save')}</Button>
        </Modal.Actions>
      </Modal>
    </>
  );
};

export default PlanManagement;