import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Button,
  Form,
  Label,
  Modal,
  Table,
  Message,
} from 'semantic-ui-react';
import {
  API,
  showError,
  showSuccess,
  timestamp2string,
} from '../helpers';

const NodeManagement = () => {
  const { t } = useTranslation();
  const [nodes, setNodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showEditModal, setShowEditModal] = useState(false);
  const [pingingNodeId, setPingingNodeId] = useState(null);
  const [editForm, setEditForm] = useState({
    node_id: 0,
    node_name: '',
    address: '',
    secret: '',
    status: 1,
  });
  const [includeDisabled, setIncludeDisabled] = useState(true);

  const loadNodes = async () => {
    setLoading(true);
    const res = await API.get('/api/cluster_node/');
    const { success, message, data } = res.data;
    if (success) {
      setNodes(data || []);
    } else {
      showError(message);
    }
    setLoading(false);
  };

  useEffect(() => {
    loadNodes();
  }, []);

  const saveNode = async () => {
    const isEdit = editForm.id > 0;
    const payload = {
      node_id: editForm.node_id,
      node_name: editForm.node_name,
      address: editForm.address,
      secret: editForm.secret,
    };
    const res = isEdit
      ? await API.put('/api/cluster_node/', payload)
      : await API.post('/api/cluster_node/', payload);
    const { success, message } = res.data;
    if (success) {
      showSuccess(isEdit ? t('node.messages.update_success') : t('node.messages.add_success'));
      setShowEditModal(false);
      loadNodes();
    } else {
      showError(message);
    }
  };

  const disableNode = async (nodeId) => {
    const res = await API.delete(`/api/cluster_node/${nodeId}/`);
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('node.messages.disable_success'));
      loadNodes();
    } else {
      showError(message);
    }
  };

  const enableNode = async (nodeId) => {
    const res = await API.post(`/api/cluster_node/${nodeId}/enable`);
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('node.messages.enable_success'));
      loadNodes();
    } else {
      showError(message);
    }
  };

  const pingNode = async (nodeId) => {
    setPingingNodeId(nodeId);
    const res = await API.get(`/api/cluster_node/ping/${nodeId}`);
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('node.messages.ping_success'));
      loadNodes();
    } else {
      showError(message || t('node.messages.ping_failed'));
    }
    setPingingNodeId(null);
  };

  const openEditModal = (node) => {
    if (node) {
      setEditForm({
        id: node.id,
        node_id: node.node_id,
        node_name: node.node_name,
        address: node.address,
        secret: node.secret_key || '',
        status: node.status,
      });
    } else {
      setEditForm({
        id: 0,
        node_id: 0,
        node_name: '',
        address: '',
        secret: '',
        status: 1,
      });
    }
    setShowEditModal(true);
  };

  const renderStatus = (node) => {
    if (node.disabled) {
      return <Label basic color='grey'>{t('node.status.disabled')}</Label>;
    }
    if (node.status === 1) {
      return <Label basic color='green'>{t('node.status.alive')}</Label>;
    }
    return <Label basic color='red'>{t('node.status.dead')}</Label>;
  };

  const renderNodeName = (node) => {
    if (node.is_self) {
      return (
        <>
          {node.node_name}
          <Label basic color='blue' style={{ marginLeft: 8 }}>
            {t('node.self')}
          </Label>
        </>
      );
    }
    return node.node_name;
  };

  const formatTime = (timestamp) => {
    if (!timestamp || timestamp === 0) return '-';
    return timestamp2string(timestamp);
  };

  const visibleNodes = includeDisabled
    ? nodes
    : nodes.filter((n) => !n.disabled);

  return (
    <>
      <Message info>
        <Message.Header>{t('node.help.title')}</Message.Header>
        <p>{t('node.help.content')}</p>
        <p style={{ marginTop: 8, color: '#666' }}>{t('node.help.warning')}</p>
      </Message>

      <Button size='small' onClick={() => openEditModal(null)}>
        {t('node.buttons.add')}
      </Button>
      <Button size='small' onClick={loadNodes} loading={loading}>
        {t('node.buttons.refresh')}
      </Button>
      <Button
        size='small'
        toggle
        active={includeDisabled}
        onClick={() => setIncludeDisabled(!includeDisabled)}
      >
        {t('node.buttons.show_disabled')}
      </Button>

      <Table basic='very' compact size='small' style={{ marginTop: 10 }}>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>ID</Table.HeaderCell>
            <Table.HeaderCell>{t('node.table.node_id')}</Table.HeaderCell>
            <Table.HeaderCell>{t('node.table.node_name')}</Table.HeaderCell>
            <Table.HeaderCell>{t('node.table.address')}</Table.HeaderCell>
            <Table.HeaderCell>{t('node.table.status')}</Table.HeaderCell>
            <Table.HeaderCell>{t('node.table.ping_failures')}</Table.HeaderCell>
            <Table.HeaderCell>{t('node.table.last_heartbeat')}</Table.HeaderCell>
            <Table.HeaderCell>{t('node.table.last_ping_attempt')}</Table.HeaderCell>
            <Table.HeaderCell>{t('node.table.actions')}</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {visibleNodes.map((node) => (
            <Table.Row key={node.id}>
              <Table.Cell>{node.id}</Table.Cell>
              <Table.Cell>{node.node_id}</Table.Cell>
              <Table.Cell>{renderNodeName(node)}</Table.Cell>
              <Table.Cell>{node.address}</Table.Cell>
              <Table.Cell>{renderStatus(node)}</Table.Cell>
              <Table.Cell>{node.ping_failures}</Table.Cell>
              <Table.Cell>{formatTime(node.last_heartbeat)}</Table.Cell>
              <Table.Cell>{formatTime(node.last_ping_attempt)}</Table.Cell>
              <Table.Cell>
                <Button
                  size='tiny'
                  onClick={() => pingNode(node.node_id)}
                  loading={pingingNodeId === node.node_id}
                  disabled={pingingNodeId !== null || node.is_self || node.disabled}
                >
                  {t('node.buttons.ping')}
                </Button>
                <Button size='tiny' onClick={() => openEditModal(node)}>
                  {t('node.buttons.edit')}
                </Button>
                {node.disabled ? (
                  <Button
                    size='tiny'
                    positive
                    onClick={() => enableNode(node.node_id)}
                  >
                    {t('node.buttons.enable')}
                  </Button>
                ) : (
                  <Modal
                    trigger={
                      <Button size='tiny' negative disabled={node.is_self}>
                        {t('node.buttons.disable')}
                      </Button>
                    }
                    header={t('node.modals.disable_confirm')}
                    content={t('node.modals.disable_content', { name: node.node_name || node.address })}
                    actions={[
                      { key: 'cancel', content: t('node.buttons.cancel'), positive: false },
                      { key: 'disable', content: t('node.buttons.confirm_disable'), positive: true, onClick: () => disableNode(node.node_id) },
                    ]}
                  />
                )}
              </Table.Cell>
            </Table.Row>
          ))}
          {visibleNodes.length === 0 && !loading && (
            <Table.Row>
              <Table.Cell colSpan={9} textAlign='center'>
                {t('node.table.empty')}
              </Table.Cell>
            </Table.Row>
          )}
        </Table.Body>
      </Table>

      <Modal open={showEditModal} onClose={() => setShowEditModal(false)} size='small'>
        <Modal.Header>{editForm.id ? t('node.modals.edit_title') : t('node.modals.add_title')}</Modal.Header>
        <Modal.Content>
          <Form>
            <Form.Input
              label={t('node.form.node_id')}
              type='number'
              value={editForm.node_id}
              onChange={(e, { value }) => setEditForm({ ...editForm, node_id: parseInt(value) || 0 })}
              disabled={editForm.id > 0}
              placeholder='1-49'
            />
            <Form.Input
              label={t('node.form.node_name')}
              value={editForm.node_name}
              onChange={(e, { value }) => setEditForm({ ...editForm, node_name: value })}
              placeholder={t('node.form.node_name_placeholder')}
            />
            <Form.Input
              label={t('node.form.address')}
              value={editForm.address}
              onChange={(e, { value }) => setEditForm({ ...editForm, address: value })}
              placeholder={t('node.form.address_placeholder')}
            />
            <Form.Input
              label={t('node.form.secret')}
              value={editForm.secret}
              onChange={(e, { value }) => setEditForm({ ...editForm, secret: value })}
              placeholder={t('node.form.secret_placeholder')}
              type='text'
            />
            {editForm.id > 0 && (
              <Message info size='small'>
                {t('node.form.secret_hint')}
              </Message>
            )}
          </Form>
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={() => setShowEditModal(false)}>{t('node.buttons.cancel')}</Button>
          <Button positive onClick={saveNode}>{t('node.buttons.save')}</Button>
        </Modal.Actions>
      </Modal>
    </>
  );
};

export default NodeManagement;