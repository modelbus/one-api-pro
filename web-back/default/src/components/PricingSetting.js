import React, { useEffect, useState } from 'react';
import { Divider, Form, Grid, Header, Button, Table, Modal, Message, Icon, Input } from 'semantic-ui-react';
import {
  API,
  showError,
  showSuccess,
} from '../helpers';

const PricingSetting = () => {
  const [modelPrices, setModelPrices] = useState([]);
  const [groupPrices, setGroupPrices] = useState([]);
  const [modelPriceModal, setModelPriceModal] = useState({ open: false, edit: null });
  const [groupPriceModal, setGroupPriceModal] = useState({ open: false, edit: null });
  const [modelPriceForm, setModelPriceForm] = useState({
    model_name: '', input_price: 0, output_price: 0, cached_price: 0, per_request_price: 0, billing_type: 'token', enabled: true,
  });
  const [groupPriceForm, setGroupPriceForm] = useState({
    group_name: '', model_name: '', discount: 1.0,
  });
  const [searchText, setSearchText] = useState('');

  const getModelPrices = async () => {
    const res = await API.get('/api/model_price/');
    const { success, message, data } = res.data;
    if (success) {
      setModelPrices(data || []);
    } else {
      showError(message);
    }
  };

  const getGroupPrices = async () => {
    const res = await API.get('/api/group_price/');
    const { success, message, data } = res.data;
    if (success) {
      setGroupPrices(data || []);
    } else {
      showError(message);
    }
  };

  useEffect(() => {
    getModelPrices().then();
    getGroupPrices().then();
  }, []);

  const saveModelPrice = async () => {
    const form = { ...modelPriceForm };
    if (!form.model_name) {
      showError('模型名称不能为空');
      return;
    }
    let res;
    if (modelPriceModal.edit) {
      form.id = modelPriceModal.edit.id;
      res = await API.put('/api/model_price/', form);
    } else {
      res = await API.post('/api/model_price/', form);
    }
    const { success, message } = res.data;
    if (success) {
      showSuccess('保存成功');
      setModelPriceModal({ open: false, edit: null });
      setModelPriceForm({ model_name: '', input_price: 0, output_price: 0, cached_price: 0, per_request_price: 0, billing_type: 'token', enabled: true });
      getModelPrices();
    } else {
      showError(message);
    }
  };

  const deleteModelPrice = async (id) => {
    const res = await API.delete(`/api/model_price/${id}`);
    const { success, message } = res.data;
    if (success) {
      showSuccess('删除成功');
      getModelPrices();
    } else {
      showError(message);
    }
  };

  const saveGroupPrice = async () => {
    const form = { ...groupPriceForm };
    if (!form.group_name) {
      showError('分组名称不能为空');
      return;
    }
    let res;
    if (groupPriceModal.edit) {
      form.id = groupPriceModal.edit.id;
      res = await API.put('/api/group_price/', form);
    } else {
      res = await API.post('/api/group_price/', form);
    }
    const { success, message } = res.data;
    if (success) {
      showSuccess('保存成功');
      setGroupPriceModal({ open: false, edit: null });
      setGroupPriceForm({ group_name: '', model_name: '', discount: 1.0 });
      getGroupPrices();
    } else {
      showError(message);
    }
  };

  const deleteGroupPrice = async (id) => {
    const res = await API.delete(`/api/group_price/${id}`);
    const { success, message } = res.data;
    if (success) {
      showSuccess('删除成功');
      getGroupPrices();
    } else {
      showError(message);
    }
  };

  const filteredPrices = searchText
    ? modelPrices.filter((p) => p.model_name.toLowerCase().includes(searchText.toLowerCase()))
    : modelPrices;

  return (
    <Grid columns={1}>
      <Grid.Column>
        <Header as='h3'>模型定价</Header>
        <Message info>
          <p>价格单位为 <strong>¥/百万tokens</strong>。计费类型为 token 的模型按实际用量计费，per_request 的模型按次计费。</p>
          <p>缓存价格为 0 时表示该模型不支持缓存折扣。</p>
        </Message>
        <Button primary size='small' onClick={() => {
          setModelPriceForm({ model_name: '', input_price: 0, output_price: 0, cached_price: 0, per_request_price: 0, billing_type: 'token', enabled: true });
          setModelPriceModal({ open: true, edit: null });
        }}>
          <Icon name='plus' /> 添加模型定价
        </Button>
        <Input
          icon='search'
          placeholder='搜索模型名称...'
          value={searchText}
          onChange={(_, { value }) => setSearchText(value)}
          size='small'
          style={{ marginLeft: '1em' }}
        />
        <Table celled compact size='small' style={{ marginTop: '1em' }}>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>模型名称</Table.HeaderCell>
              <Table.HeaderCell>输入价格</Table.HeaderCell>
              <Table.HeaderCell>输出价格</Table.HeaderCell>
              <Table.HeaderCell>缓存价格</Table.HeaderCell>
              <Table.HeaderCell>按次价格</Table.HeaderCell>
              <Table.HeaderCell>计费类型</Table.HeaderCell>
              <Table.HeaderCell>启用</Table.HeaderCell>
              <Table.HeaderCell>操作</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {filteredPrices.map((p) => (
              <Table.Row key={p.id} negative={!p.enabled}>
                <Table.Cell>{p.model_name}</Table.Cell>
                <Table.Cell>{p.input_price}</Table.Cell>
                <Table.Cell>{p.output_price}</Table.Cell>
                <Table.Cell>{p.cached_price}</Table.Cell>
                <Table.Cell>{p.per_request_price}</Table.Cell>
                <Table.Cell>{p.billing_type === 'per_request' ? '按次' : '按Token'}</Table.Cell>
                <Table.Cell>{p.enabled ? '是' : '否'}</Table.Cell>
                <Table.Cell>
                  <Button size='mini' primary onClick={() => {
                    setModelPriceForm({
                      model_name: p.model_name,
                      input_price: p.input_price,
                      output_price: p.output_price,
                      cached_price: p.cached_price,
                      per_request_price: p.per_request_price,
                      billing_type: p.billing_type,
                      enabled: p.enabled,
                    });
                    setModelPriceModal({ open: true, edit: p });
                  }}>编辑</Button>
                  <Button size='mini' negative onClick={() => deleteModelPrice(p.id)}>删除</Button>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>

        <Modal open={modelPriceModal.open} onClose={() => setModelPriceModal({ open: false, edit: null })} size='small'>
          <Modal.Header>{modelPriceModal.edit ? '编辑模型定价' : '添加模型定价'}</Modal.Header>
          <Modal.Content>
            <Form>
              <Form.Input label='模型名称' name='model_name' value={modelPriceForm.model_name}
                onChange={(_, { value }) => setModelPriceForm(f => ({ ...f, model_name: value }))}
                placeholder='例如 gpt-4o' disabled={!!modelPriceModal.edit} />
              <Form.Group widths='equal'>
                <Form.Input label='输入价格 (¥/百万tokens)' name='input_price' value={modelPriceForm.input_price}
                  onChange={(_, { value }) => setModelPriceForm(f => ({ ...f, input_price: parseFloat(value) || 0 }))}
                  type='number' step='0.001' min='0' />
                <Form.Input label='输出价格 (¥/百万tokens)' name='output_price' value={modelPriceForm.output_price}
                  onChange={(_, { value }) => setModelPriceForm(f => ({ ...f, output_price: parseFloat(value) || 0 }))}
                  type='number' step='0.001' min='0' />
                <Form.Input label='缓存价格 (¥/百万tokens)' name='cached_price' value={modelPriceForm.cached_price}
                  onChange={(_, { value }) => setModelPriceForm(f => ({ ...f, cached_price: parseFloat(value) || 0 }))}
                  type='number' step='0.001' min='0' />
              </Form.Group>
              <Form.Group widths='equal'>
                <Form.Input label='按次价格 (¥/次)' name='per_request_price' value={modelPriceForm.per_request_price}
                  onChange={(_, { value }) => setModelPriceForm(f => ({ ...f, per_request_price: parseFloat(value) || 0 }))}
                  type='number' step='0.001' min='0' />
                <Form.Select label='计费类型' value={modelPriceForm.billing_type}
                  options={[
                    { key: 'token', text: '按Token', value: 'token' },
                    { key: 'per_request', text: '按次', value: 'per_request' },
                  ]}
                  onChange={(_, { value }) => setModelPriceForm(f => ({ ...f, billing_type: value }))} />
                <Form.Select label='启用' value={modelPriceForm.enabled ? 'true' : 'false'}
                  options={[
                    { key: 'true', text: '是', value: 'true' },
                    { key: 'false', text: '否', value: 'false' },
                  ]}
                  onChange={(_, { value }) => setModelPriceForm(f => ({ ...f, enabled: value === 'true' }))} />
              </Form.Group>
            </Form>
          </Modal.Content>
          <Modal.Actions>
            <Button onClick={() => setModelPriceModal({ open: false, edit: null })}>取消</Button>
            <Button primary onClick={saveModelPrice}>保存</Button>
          </Modal.Actions>
        </Modal>

        <Divider />
        <Header as='h3'>分组折扣</Header>
        <Message info>
          <p>折扣系数为 1.0 表示无折扣，0.8 表示八折。模型名称留空表示该分组的默认折扣。</p>
        </Message>
        <Button primary size='small' onClick={() => {
          setGroupPriceForm({ group_name: '', model_name: '', discount: 1.0 });
          setGroupPriceModal({ open: true, edit: null });
        }}>
          <Icon name='plus' /> 添加分组折扣
        </Button>
        <Table celled compact size='small' style={{ marginTop: '1em' }}>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>分组名称</Table.HeaderCell>
              <Table.HeaderCell>模型名称</Table.HeaderCell>
              <Table.HeaderCell>折扣系数</Table.HeaderCell>
              <Table.HeaderCell>操作</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {groupPrices.map((p) => (
              <Table.Row key={p.id}>
                <Table.Cell>{p.group_name}</Table.Cell>
                <Table.Cell>{p.model_name || '(默认)'}</Table.Cell>
                <Table.Cell>{p.discount}</Table.Cell>
                <Table.Cell>
                  <Button size='mini' primary onClick={() => {
                    setGroupPriceForm({
                      group_name: p.group_name,
                      model_name: p.model_name,
                      discount: p.discount,
                    });
                    setGroupPriceModal({ open: true, edit: p });
                  }}>编辑</Button>
                  <Button size='mini' negative onClick={() => deleteGroupPrice(p.id)}>删除</Button>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>

        <Modal open={groupPriceModal.open} onClose={() => setGroupPriceModal({ open: false, edit: null })} size='small'>
          <Modal.Header>{groupPriceModal.edit ? '编辑分组折扣' : '添加分组折扣'}</Modal.Header>
          <Modal.Content>
            <Form>
              <Form.Group widths='equal'>
                <Form.Input label='分组名称' name='group_name' value={groupPriceForm.group_name}
                  onChange={(_, { value }) => setGroupPriceForm(f => ({ ...f, group_name: value }))}
                  placeholder='例如 default, vip' />
                <Form.Input label='模型名称 (留空为默认折扣)' name='model_name' value={groupPriceForm.model_name}
                  onChange={(_, { value }) => setGroupPriceForm(f => ({ ...f, model_name: value }))}
                  placeholder='留空表示该分组所有模型的默认折扣' />
                <Form.Input label='折扣系数' name='discount' value={groupPriceForm.discount}
                  onChange={(_, { value }) => setGroupPriceForm(f => ({ ...f, discount: parseFloat(value) || 1.0 }))}
                  type='number' step='0.01' min='0' max='1' />
              </Form.Group>
            </Form>
          </Modal.Content>
          <Modal.Actions>
            <Button onClick={() => setGroupPriceModal({ open: false, edit: null })}>取消</Button>
            <Button primary onClick={saveGroupPrice}>保存</Button>
          </Modal.Actions>
        </Modal>
      </Grid.Column>
    </Grid>
  );
};

export default PricingSetting;