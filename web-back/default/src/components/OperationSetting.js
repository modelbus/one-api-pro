import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Divider, Form, Grid, Header, Button } from 'semantic-ui-react';
import {
  API,
  showError,
  showSuccess,
  timestamp2string,
} from '../helpers';

const OperationSetting = () => {
  const { t } = useTranslation();
  let now = new Date();
  let [inputs, setInputs] = useState({
    QuotaForNewUser: 0,
    QuotaForInviter: 0,
    QuotaForInvitee: 0,
    QuotaRemindThreshold: 0,
    PreConsumedQuota: 0,
    TopUpLink: '',
    ChatLink: '',
    QuotaPerUnit: 0,
    AutomaticDisableChannelEnabled: '',
    AutomaticEnableChannelEnabled: '',
    ChannelDisableThreshold: 0,
    LogConsumeEnabled: '',
    DisplayInCurrencyEnabled: '',
    DisplayTokenStatEnabled: '',
    ApproximateTokenEnabled: '',
    RetryTimes: 0,
    ChannelDefaultCooldownSeconds: 60,
    ChannelMaxCooldownSeconds: 600,
    ChannelConcurrencyEnabled: '',
    ChannelStickySessionEnabled: '',
  });
  const [errorNext, setErrorNext] = useState({
    passthrough: true,
    retry: true,
    disable: true,
    cooldown: true,
  });
  const [originInputs, setOriginInputs] = useState({});
  let [loading, setLoading] = useState(false);
  let [historyTimestamp, setHistoryTimestamp] = useState(
    timestamp2string(now.getTime() / 1000 - 30 * 24 * 3600)
  );

  const getOptions = async () => {
    const res = await API.get('/api/option/');
    const { success, message, data } = res.data;
    if (success) {
      let newInputs = {};
      data.forEach((item) => {
        if (item.key === 'ErrorNext') {
          try {
            const parsed = JSON.parse(item.value);
            setErrorNext(parsed);
          } catch (e) {}
          return;
        }
        if (item.value === '{}') {
          item.value = '';
        }
        newInputs[item.key] = item.value;
      });
      setInputs(newInputs);
      setOriginInputs(newInputs);
    } else {
      showError(message);
    }
  };

  useEffect(() => {
    getOptions().then();
  }, []);

  const updateOption = async (key, value) => {
    setLoading(true);
    if (key.endsWith('Enabled')) {
      value = inputs[key] === 'true' ? 'false' : 'true';
    }
    const res = await API.put('/api/option/', {
      key,
      value,
    });
    const { success, message } = res.data;
    if (success) {
      setInputs((inputs) => ({ ...inputs, [key]: value }));
    } else {
      showError(message);
    }
    setLoading(false);
  };

  const handleInputChange = async (e, { name, value }) => {
    if (name.endsWith('Enabled')) {
      await updateOption(name, value);
    } else {
      setInputs((inputs) => ({ ...inputs, [name]: value }));
    }
  };

  const submitConfig = async (group) => {
    switch (group) {
      case 'monitor':
        if (originInputs['ChannelDisableThreshold'] !== inputs.ChannelDisableThreshold) {
          await updateOption('ChannelDisableThreshold', inputs.ChannelDisableThreshold);
        }
        if (originInputs['QuotaRemindThreshold'] !== inputs.QuotaRemindThreshold) {
          await updateOption('QuotaRemindThreshold', inputs.QuotaRemindThreshold);
        }
        break;
      case 'quota':
        if (originInputs['QuotaForNewUser'] !== inputs.QuotaForNewUser) {
          await updateOption('QuotaForNewUser', inputs.QuotaForNewUser);
        }
        if (originInputs['QuotaForInvitee'] !== inputs.QuotaForInvitee) {
          await updateOption('QuotaForInvitee', inputs.QuotaForInvitee);
        }
        if (originInputs['QuotaForInviter'] !== inputs.QuotaForInviter) {
          await updateOption('QuotaForInviter', inputs.QuotaForInviter);
        }
        if (originInputs['PreConsumedQuota'] !== inputs.PreConsumedQuota) {
          await updateOption('PreConsumedQuota', inputs.PreConsumedQuota);
        }
        break;
      case 'general':
        if (originInputs['TopUpLink'] !== inputs.TopUpLink) {
          await updateOption('TopUpLink', inputs.TopUpLink);
        }
        if (originInputs['ChatLink'] !== inputs.ChatLink) {
          await updateOption('ChatLink', inputs.ChatLink);
        }
        if (originInputs['QuotaPerUnit'] !== inputs.QuotaPerUnit) {
          await updateOption('QuotaPerUnit', inputs.QuotaPerUnit);
        }
        if (originInputs['RetryTimes'] !== inputs.RetryTimes) {
          await updateOption('RetryTimes', inputs.RetryTimes);
        }
        if (originInputs['ChannelDefaultCooldownSeconds'] !== inputs.ChannelDefaultCooldownSeconds) {
          await updateOption('ChannelDefaultCooldownSeconds', inputs.ChannelDefaultCooldownSeconds);
        }
        if (originInputs['ChannelMaxCooldownSeconds'] !== inputs.ChannelMaxCooldownSeconds) {
          await updateOption('ChannelMaxCooldownSeconds', inputs.ChannelMaxCooldownSeconds);
        }
        if (originInputs['ChannelConcurrencyEnabled'] !== inputs.ChannelConcurrencyEnabled) {
          await updateOption('ChannelConcurrencyEnabled', inputs.ChannelConcurrencyEnabled);
        }
        if (originInputs['ChannelStickySessionEnabled'] !== inputs.ChannelStickySessionEnabled) {
          await updateOption('ChannelStickySessionEnabled', inputs.ChannelStickySessionEnabled);
        }
        await updateOption('ErrorNext', JSON.stringify(errorNext));
        break;
    }
  };

  const deleteHistoryLogs = async () => {
    const res = await API.delete(
      `/api/log/?target_timestamp=${Date.parse(historyTimestamp) / 1000}`
    );
    const { success, message, data } = res.data;
    if (success) {
      showSuccess(`${data} 条日志已清理！`);
      return;
    }
    showError('日志清理失败：' + message);
  };

  return (
    <Grid columns={1}>
      <Grid.Column>
        <Form loading={loading}>
          <Header as='h3'>{t('setting.operation.quota.title')}</Header>
          <Form.Group widths='equal'>
            <Form.Input
              label={t('setting.operation.quota.new_user')}
              name='QuotaForNewUser'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.QuotaForNewUser}
              type='number'
              min='0'
              placeholder={t('setting.operation.quota.new_user_placeholder')}
            />
            <Form.Input
              label={t('setting.operation.quota.pre_consume')}
              name='PreConsumedQuota'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.PreConsumedQuota}
              type='number'
              min='0'
              placeholder={t('setting.operation.quota.pre_consume_placeholder')}
            />
            <Form.Input
              label={t('setting.operation.quota.inviter_reward')}
              name='QuotaForInviter'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.QuotaForInviter}
              type='number'
              min='0'
              placeholder={t('setting.operation.quota.inviter_reward_placeholder')}
            />
            <Form.Input
              label={t('setting.operation.quota.invitee_reward')}
              name='QuotaForInvitee'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.QuotaForInvitee}
              type='number'
              min='0'
              placeholder={t('setting.operation.quota.invitee_reward_placeholder')}
            />
          </Form.Group>
          <Form.Button onClick={() => { submitConfig('quota').then(); }}>
            {t('setting.operation.quota.buttons.save')}
          </Form.Button>

          <Divider />
          <Header as='h3'>{t('setting.operation.log.title')}</Header>
          <Form.Group inline>
            <Form.Checkbox
              checked={inputs.LogConsumeEnabled === 'true'}
              label={t('setting.operation.log.enable_consume')}
              name='LogConsumeEnabled'
              onChange={handleInputChange}
            />
          </Form.Group>
          <Form.Group widths={4}>
            <Form.Input
              label={t('setting.operation.log.target_time')}
              value={historyTimestamp}
              type='datetime-local'
              name='history_timestamp'
              onChange={(e, { name, value }) => { setHistoryTimestamp(value); }}
            />
          </Form.Group>
          <Form.Button onClick={() => { deleteHistoryLogs().then(); }}>
            {t('setting.operation.log.buttons.clean')}
          </Form.Button>

          <Divider />
          <Header as='h3'>{t('setting.operation.monitor.title')}</Header>
          <Form.Group widths={3}>
            <Form.Input
              label={t('setting.operation.monitor.max_response_time')}
              name='ChannelDisableThreshold'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.ChannelDisableThreshold}
              type='number'
              min='0'
              placeholder={t('setting.operation.monitor.max_response_time_placeholder')}
            />
            <Form.Input
              label={t('setting.operation.monitor.quota_reminder')}
              name='QuotaRemindThreshold'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.QuotaRemindThreshold}
              type='number'
              min='0'
              placeholder={t('setting.operation.monitor.quota_reminder_placeholder')}
            />
          </Form.Group>
          <Form.Group inline>
            <Form.Checkbox
              checked={inputs.AutomaticDisableChannelEnabled === 'true'}
              label={t('setting.operation.monitor.auto_disable')}
              name='AutomaticDisableChannelEnabled'
              onChange={handleInputChange}
            />
            <Form.Checkbox
              checked={inputs.AutomaticEnableChannelEnabled === 'true'}
              label={t('setting.operation.monitor.auto_enable')}
              name='AutomaticEnableChannelEnabled'
              onChange={handleInputChange}
            />
          </Form.Group>
          <Form.Button onClick={() => { submitConfig('monitor').then(); }}>
            {t('setting.operation.monitor.buttons.save')}
          </Form.Button>

          <Divider />
          <Header as='h3'>{t('setting.operation.general.title')}</Header>
          <Form.Group widths={4}>
            <Form.Input
              label={t('setting.operation.general.topup_link')}
              name='TopUpLink'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.TopUpLink}
              type='link'
              placeholder={t('setting.operation.general.topup_link_placeholder')}
            />
            <Form.Input
              label={t('setting.operation.general.chat_link')}
              name='ChatLink'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.ChatLink}
              type='link'
              placeholder={t('setting.operation.general.chat_link_placeholder')}
            />
            <Form.Input
              label={t('setting.operation.general.quota_per_unit')}
              name='QuotaPerUnit'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.QuotaPerUnit}
              type='number'
              step='0.01'
              placeholder={t('setting.operation.general.quota_per_unit_placeholder')}
            />
            <Form.Input
              label={t('setting.operation.general.retry_times')}
              name='RetryTimes'
              type={'number'}
              step='1'
              min='0'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.RetryTimes}
              placeholder={t('setting.operation.general.retry_times_placeholder')}
            />
          </Form.Group>
          <Form.Group inline>
            <Form.Checkbox
              checked={inputs.DisplayInCurrencyEnabled === 'true'}
              label={t('setting.operation.general.display_in_currency')}
              name='DisplayInCurrencyEnabled'
              onChange={handleInputChange}
            />
            <Form.Checkbox
              checked={inputs.DisplayTokenStatEnabled === 'true'}
              label={t('setting.operation.general.display_token_stat')}
              name='DisplayTokenStatEnabled'
              onChange={handleInputChange}
            />
            <Form.Checkbox
              checked={inputs.ApproximateTokenEnabled === 'true'}
              label={t('setting.operation.general.approximate_token')}
              name='ApproximateTokenEnabled'
              onChange={handleInputChange}
            />
          </Form.Group>
          <Divider />
          <Header as='h3'>渠道路由设置</Header>
          <Form.Group>
            <Form.Input
              label='默认冷却时间（秒）'
              name='ChannelDefaultCooldownSeconds'
              type={'number'}
              step='1'
              min='0'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.ChannelDefaultCooldownSeconds}
              placeholder='收到429/529后渠道冷却时间，默认60秒'
            />
            <Form.Input
              label='最大冷却时间（秒）'
              name='ChannelMaxCooldownSeconds'
              type={'number'}
              step='1'
              min='0'
              onChange={handleInputChange}
              autoComplete='new-password'
              value={inputs.ChannelMaxCooldownSeconds}
              placeholder='冷却时间上限，默认600秒'
            />
          </Form.Group>
          <Form.Group inline>
            <Form.Checkbox
              checked={inputs.ChannelConcurrencyEnabled === 'true'}
              label='启用渠道并发限制'
              name='ChannelConcurrencyEnabled'
              onChange={handleInputChange}
            />
            <Form.Checkbox
              checked={inputs.ChannelStickySessionEnabled === 'true'}
              label='启用粘性会话（同一用户+模型路由到同一渠道）'
              name='ChannelStickySessionEnabled'
              onChange={handleInputChange}
            />
          </Form.Group>
          <Header as='h4'>错误响应策略（ErrorNext）</Header>
          <p style={{color:'#888', fontSize:'0.9em'}}>控制上游返回错误时系统的行为。关闭后该类错误将直接返回给客户端，不执行对应策略。</p>
          <Form.Group inline>
            <Form.Checkbox
              checked={errorNext.passthrough === true}
              label='直通 (400/404/422)'
              onChange={() => setErrorNext(e => ({...e, passthrough: !e.passthrough}))}
            />
            <Form.Checkbox
              checked={errorNext.retry === true}
              label='重试 (500/502/503)'
              onChange={() => setErrorNext(e => ({...e, retry: !e.retry}))}
            />
            <Form.Checkbox
              checked={errorNext.disable === true}
              label='禁用渠道 (401/402/403)'
              onChange={() => setErrorNext(e => ({...e, disable: !e.disable}))}
            />
            <Form.Checkbox
              checked={errorNext.cooldown === true}
              label='冷却+重试 (429/529)'
              onChange={() => setErrorNext(e => ({...e, cooldown: !e.cooldown}))}
            />
          </Form.Group>
          <Form.Button onClick={() => { submitConfig('general').then(); }}>
            {t('setting.operation.general.buttons.save')}
          </Form.Button>
        </Form>
      </Grid.Column>
    </Grid>
  );
};

export default OperationSetting;