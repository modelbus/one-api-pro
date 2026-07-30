import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Button,
  Form,
  Label,
  Modal,
  Pagination,
  Popup,
  Progress,
  Table,
} from 'semantic-ui-react';
import {
  API,
  isAdmin,
  showError,
  showSuccess,
  timestamp2string,
} from '../helpers';
import { ITEMS_PER_PAGE } from '../constants';

const MODEL_COLORS = ['#2196F3', '#4CAF50', '#FF9800', '#E91E63', '#9C27B0', '#00BCD4', '#FF5722', '#795548'];

function renderStatus(status, t) {
  switch (status) {
    case 1:
      return <Label basic color='green'>{t('subscription.status.active')}</Label>;
    case 0:
      return <Label basic color='grey'>{t('subscription.status.expired')}</Label>;
    default:
      return <Label basic color='black'>{t('subscription.status.unknown')}</Label>;
  }
}

function renderBillingType(type, t) {
  switch (type) {
    case 'request':
      return t('subscription.billing_type.request');
    case 'token':
      return t('subscription.billing_type.token');
    default:
      return type;
  }
}

function renderWindowType(wt, t) {
  switch (wt) {
    case 'period':
      return t('subscription.usage.period');
    case 'week':
      return t('subscription.usage.week');
    case 'month':
      return t('subscription.usage.month');
    default:
      return wt;
  }
}

function formatPercent(val) {
  if (val === null || val === undefined) return '-';
  return val.toFixed(2) + '%';
}

function formatNumber(num) {
  if (num === null || num === undefined) return '-';
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
}

function MultiSegmentProgress({ segments, total }) {
  if (!segments || segments.length === 0) return null;
  return (
    <div style={{ position: 'relative', height: 24, backgroundColor: '#e8e8e8', borderRadius: 4, overflow: 'hidden' }}>
      {segments.map((seg, i) => {
        const pct = total > 0 ? (seg.value / total) * 100 : 0;
        return (
          <div
            key={i}
            title={`${seg.label}: ${formatNumber(seg.value)} (${pct.toFixed(2)}%)`}
            style={{
              position: 'absolute',
              left: `${segments.slice(0, i).reduce((a, s) => a + (total > 0 ? (s.value / total) * 100 : 0), 0)}%`,
              width: `${pct}%`,
              height: '100%',
              backgroundColor: seg.color,
              opacity: 0.85,
              transition: 'width 0.3s ease',
            }}
          />
        );
      })}
    </div>
  );
}

const SubscriptionsTable = () => {
  const { t } = useTranslation();
  const isAdminUser = isAdmin();
  const [subscriptions, setSubscriptions] = useState([]);
  const [plans, setPlans] = useState([]);
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activePage, setActivePage] = useState(1);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [searching, setSearching] = useState(false);
  const [showAddModal, setShowAddModal] = useState(false);
  const [addForm, setAddForm] = useState({ user_id: '', plan_id: '', billing_type: 'token', duration_days: '', notes: '' });
  const [showRenewModal, setShowRenewModal] = useState(false);
  const [renewForm, setRenewForm] = useState({ id: '', end_time: '', notes: '' });
  const [showUsageModal, setShowUsageModal] = useState(false);
  const [usageData, setUsageData] = useState(null);
  const [expandedWindows, setExpandedWindows] = useState({});

  const loadSubscriptions = async (startIdx) => {
    let url;
    if (isAdminUser) {
      url = `/api/subscription/?p=${startIdx}`;
    } else {
      url = '/api/subscription/self';
    }
    const res = await API.get(url);
    const { success, message, data } = res.data;
    if (success) {
      if (isAdminUser) {
        if (startIdx === 0) {
          setSubscriptions(data);
        } else {
          setSubscriptions([...subscriptions, ...data]);
        }
      } else {
        setSubscriptions(Array.isArray(data) ? data : []);
      }
    } else {
      showError(message);
    }
    setLoading(false);
  };

  const loadPlans = async () => {
    const res = await API.get('/api/plan/?p=0');
    const { success, data } = res.data;
    if (success) {
      setPlans(data);
    }
  };

  const loadUsers = async () => {
    const res = await API.get('/api/user/?p=0');
    const { success, data } = res.data;
    if (success) {
      setUsers(data);
    }
  };

  const onPaginationChange = (e, { activePage }) => {
    (async () => {
      if (isAdminUser && activePage === Math.ceil(subscriptions.length / ITEMS_PER_PAGE) + 1) {
        await loadSubscriptions(activePage - 1);
      }
      setActivePage(activePage);
    })();
  };

  useEffect(() => {
    loadSubscriptions(0).then().catch((reason) => { showError(reason); });
    if (isAdminUser) {
      loadPlans();
      loadUsers();
    }
  }, []);

  const searchSubscriptions = async () => {
    if (!isAdminUser) return;
    if (searchKeyword === '') {
      await loadSubscriptions(0);
      setActivePage(1);
      return;
    }
    setSearching(true);
    const res = await API.get(`/api/subscription/search?keyword=${searchKeyword}`);
    const { success, message, data } = res.data;
    if (success) {
      setSubscriptions(data);
      setActivePage(1);
    } else {
      showError(message);
    }
    setSearching(false);
  };

  const addSubscription = async () => {
    const res = await API.post('/api/subscription/', {
      ...addForm,
      user_id: parseInt(addForm.user_id),
      plan_id: parseInt(addForm.plan_id),
      duration_days: addForm.duration_days ? parseInt(addForm.duration_days) : 0,
    });
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('subscription.messages.add_success'));
      setShowAddModal(false);
      setAddForm({ user_id: '', plan_id: '', billing_type: 'token', duration_days: '', notes: '' });
      loadSubscriptions(0);
    } else {
      showError(message);
    }
  };

  const renewSubscription = async () => {
    const res = await API.put('/api/subscription/', renewForm);
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('subscription.messages.renew_success'));
      setShowRenewModal(false);
      loadSubscriptions(0);
    } else {
      showError(message);
    }
  };

  const expireSubscription = async (id) => {
    const res = await API.put('/api/subscription/', { id, status: 0 });
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('subscription.messages.expire_success'));
      loadSubscriptions(0);
    } else {
      showError(message);
    }
  };

  const deleteSubscription = async (id) => {
    const res = await API.delete(`/api/subscription/${id}/`);
    const { success, message } = res.data;
    if (success) {
      showSuccess(t('subscription.messages.delete_success'));
      loadSubscriptions(0);
    } else {
      showError(message);
    }
  };

  const viewUsage = async (id) => {
    const res = await API.get(`/api/subscription/${id}/usage`);
    const { success, data, message } = res.data;
    if (success) {
      setUsageData(data);
      setShowUsageModal(true);
    } else {
      showError(message);
    }
  };

  const refresh = async () => {
    setLoading(true);
    await loadSubscriptions(0);
    setActivePage(1);
  };

  const getUserName = (userId) => {
    const user = users.find(u => u.id === userId);
    return user ? user.username : userId;
  };

  const toggleWindow = (wt) => {
    setExpandedWindows(prev => ({ ...prev, [wt]: !prev[wt] }));
  };

  const renderUsageProgress = () => {
    if (!usageData || !usageData.weighted) return null;
    const { weighted, next_reset, billing_type } = usageData;
    const windowTypes = ['period', 'week', 'month'];

    return (
      <div style={{ marginBottom: 16 }}>
        {windowTypes.map(wt => {
          const wVal = weighted[wt] || 0;
          const pct = Math.min(wVal, 100);
          const isExhausted = wVal >= 100;
          const isExpanded = expandedWindows[wt];
          const resetTime = next_reset ? next_reset[wt] : null;
          const modelDetails = usageData.model_usage ? usageData.model_usage[wt] : [];

          // Build segments for the progress bar
          let segments = [];
          if (modelDetails && modelDetails.length > 0) {
            const isToken = billing_type === 'token';
            modelDetails.forEach((md, i) => {
              segments.push({
                label: md.model,
                value: isToken ? (md.prompt_tokens + md.completion_tokens) : md.requests,
                color: MODEL_COLORS[i % MODEL_COLORS.length],
              });
            });
          }

          // Determine limit for the progress bar total
          const limit = getLimitForWindow(wt, usageData.limits, billing_type);

          return (
            <div key={wt} style={{ marginBottom: 12 }}>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 4 }}>
                <span style={{ fontWeight: 'bold', cursor: 'pointer' }} onClick={() => toggleWindow(wt)}>
                  {isExpanded ? '▼' : '▶'} {renderWindowType(wt, t)}
                </span>
                <span>
                  <span style={{ fontWeight: 'bold', color: isExhausted ? 'red' : 'inherit', fontSize: '0.95em' }}>
                    {formatPercent(wVal)}
                  </span>
                  {resetTime && (
                    <span style={{ fontSize: '0.8em', color: '#888', marginLeft: 8 }}>
                      {t('subscription.usage.next_reset')}: {timestamp2string(resetTime)}
                    </span>
                  )}
                </span>
              </div>
              <Progress
                percent={pct}
                label={formatPercent(wVal)}
                error={isExhausted}
                warning={!isExhausted && pct >= 80}
                success={pct < 50}
                size='small'
                style={{ marginBottom: 2 }}
              />
              {segments.length > 0 && (
                <MultiSegmentProgress segments={segments} total={limit > 0 ? limit : 1} />
              )}
              {segments.length > 0 && (
                <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', marginTop: 4, fontSize: '0.8em' }}>
                  {segments.map((seg, i) => (
                    <span key={i} style={{ display: 'flex', alignItems: 'center', gap: 3 }}>
                      <span style={{ display: 'inline-block', width: 10, height: 10, backgroundColor: seg.color, borderRadius: 2 }} />
                      <span>{seg.label}: {formatNumber(seg.value)}</span>
                    </span>
                  ))}
                </div>
              )}
              {isExpanded && modelDetails && modelDetails.length > 0 && (
                <Table compact size='small' style={{ marginTop: 4 }}>
                  <Table.Header>
                    <Table.Row>
                      <Table.HeaderCell>{t('subscription.usage.model')}</Table.HeaderCell>
                      <Table.HeaderCell>{t('subscription.usage.requests')}</Table.HeaderCell>
                      <Table.HeaderCell>%</Table.HeaderCell>
                      <Table.HeaderCell>{t('subscription.usage.prompt_tokens')}</Table.HeaderCell>
                      <Table.HeaderCell>{t('subscription.usage.completion_tokens')}</Table.HeaderCell>
                      <Table.HeaderCell>{t('subscription.usage.cached_tokens')}</Table.HeaderCell>
                      <Table.HeaderCell>%</Table.HeaderCell>
                    </Table.Row>
                  </Table.Header>
                  <Table.Body>
                    {modelDetails.map((md, i) => (
                      <Table.Row key={i}>
                        <Table.Cell>
                          <span style={{ display: 'inline-block', width: 10, height: 10, backgroundColor: MODEL_COLORS[i % MODEL_COLORS.length], borderRadius: 2, marginRight: 4 }} />
                          {md.model}
                        </Table.Cell>
                        <Table.Cell>{formatNumber(md.requests)}</Table.Cell>
                        <Table.Cell>{md.request_percent > 0 ? formatPercent(md.request_percent) : '-'}</Table.Cell>
                        <Table.Cell>{formatNumber(md.prompt_tokens)}</Table.Cell>
                        <Table.Cell>{formatNumber(md.completion_tokens)}</Table.Cell>
                        <Table.Cell>{formatNumber(md.cached_tokens)}</Table.Cell>
                        <Table.Cell>{md.token_percent > 0 ? formatPercent(md.token_percent) : '-'}</Table.Cell>
                      </Table.Row>
                    ))}
                  </Table.Body>
                </Table>
              )}
            </div>
          );
        })}
      </div>
    );
  };

  const getLimitForWindow = (wt, limits, billingType) => {
    if (!limits) return 0;
    // Get the other (fallback) rule or first rule
    let rule = limits['other'];
    if (!rule) {
      const keys = Object.keys(limits);
      if (keys.length === 0) return 0;
      rule = limits[keys[0]];
    }
    if (billingType === 'token') {
      switch (wt) {
        case 'period': return rule.token_period;
        case 'week': return rule.token_week;
        case 'month': return rule.token_month;
      }
    } else {
      switch (wt) {
        case 'period': return rule.request_period;
        case 'week': return rule.request_week;
        case 'month': return rule.request_month;
      }
    }
    return 0;
  };

  return (
    <>
      {isAdminUser && (
        <Form onSubmit={searchSubscriptions}>
          <Form.Input
            icon='search'
            fluid
            iconPosition='left'
            placeholder={t('subscription.search')}
            value={searchKeyword}
            loading={searching}
            onChange={(e, { value }) => setSearchKeyword(value.trim())}
          />
        </Form>
      )}

      <Table basic='very' compact size='small'>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>ID</Table.HeaderCell>
            {isAdminUser && (
              <Table.HeaderCell>{t('subscription.table.user')}</Table.HeaderCell>
            )}
            <Table.HeaderCell>{t('subscription.table.plan')}</Table.HeaderCell>
            <Table.HeaderCell>{t('subscription.table.billing_type')}</Table.HeaderCell>
            <Table.HeaderCell>{t('subscription.table.start_time')}</Table.HeaderCell>
            <Table.HeaderCell>{t('subscription.table.end_time')}</Table.HeaderCell>
            <Table.HeaderCell>{t('subscription.table.status')}</Table.HeaderCell>
            <Table.HeaderCell>{t('subscription.table.actions')}</Table.HeaderCell>
          </Table.Row>
        </Table.Header>

        <Table.Body>
          {subscriptions
            .slice((activePage - 1) * ITEMS_PER_PAGE, activePage * ITEMS_PER_PAGE)
            .map((sub) => (
              <Table.Row key={sub.id}>
                <Table.Cell>{sub.id}</Table.Cell>
                {isAdminUser && (
                  <Table.Cell>{getUserName(sub.user_id)}</Table.Cell>
                )}
                <Table.Cell>{sub.plan ? sub.plan.name : sub.plan_id}</Table.Cell>
                <Table.Cell>{renderBillingType(sub.billing_type, t)}</Table.Cell>
                <Table.Cell>{timestamp2string(sub.start_time)}</Table.Cell>
                <Table.Cell>{timestamp2string(sub.end_time)}</Table.Cell>
                <Table.Cell>{renderStatus(sub.status, t)}</Table.Cell>
                <Table.Cell>
                  <Button size='tiny' onClick={() => viewUsage(sub.id)}>
                    {t('subscription.buttons.usage')}
                  </Button>
                  {isAdminUser && sub.status === 1 && (
                    <>
                      <Button size='tiny' onClick={() => {
                        setRenewForm({ id: sub.id, end_time: sub.end_time + 30 * 86400, notes: '' });
                        setShowRenewModal(true);
                      }}>
                        {t('subscription.buttons.renew')}
                      </Button>
                      <Button size='tiny' negative onClick={() => expireSubscription(sub.id)}>
                        {t('subscription.buttons.expire')}
                      </Button>
                    </>
                  )}
                  {isAdminUser && (
                    <Popup
                      trigger={<Button size='tiny' negative>{t('subscription.buttons.delete')}</Button>}
                      on='click' flowing hoverable>
                      <Button negative onClick={() => deleteSubscription(sub.id)}>
                        {t('subscription.buttons.confirm_delete')}
                      </Button>
                    </Popup>
                  )}
                </Table.Cell>
              </Table.Row>
            ))}
        </Table.Body>

        <Table.Footer>
          <Table.Row>
            <Table.HeaderCell colSpan={isAdminUser ? '8' : '7'}>
              {isAdminUser && (
                <Button size='small' onClick={() => { loadPlans(); loadUsers(); setShowAddModal(true); }}>
                  {t('subscription.buttons.add')}
                </Button>
              )}
              <Button size='small' onClick={refresh} loading={loading}>
                {t('subscription.buttons.refresh')}
              </Button>
              {isAdminUser && (
                <Pagination
                  floated='right'
                  activePage={activePage}
                  onPageChange={onPaginationChange}
                  size='small'
                  siblingRange={1}
                  totalPages={
                    Math.ceil(subscriptions.length / ITEMS_PER_PAGE) +
                    (subscriptions.length % ITEMS_PER_PAGE === 0 ? 1 : 0)
                  }
                />
              )}
            </Table.HeaderCell>
          </Table.Row>
        </Table.Footer>
      </Table>

      {/* Add Subscription Modal */}
      {isAdminUser && (
        <Modal open={showAddModal} onClose={() => setShowAddModal(false)}>
          <Modal.Header>{t('subscription.modals.add_title')}</Modal.Header>
          <Modal.Content>
            <Form>
              <Form.Select
                label={t('subscription.modals.user')}
                options={users.map(u => ({ key: u.id, text: u.username, value: u.id }))}
                value={addForm.user_id}
                onChange={(e, { value }) => setAddForm({ ...addForm, user_id: value })}
              />
              <Form.Select
                label={t('subscription.modals.plan')}
                options={plans.map(p => ({ key: p.id, text: p.name + (p.recommended ? ' *' : ''), value: p.id }))}
                value={addForm.plan_id}
                onChange={(e, { value }) => setAddForm({ ...addForm, plan_id: value })}
              />
              <Form.Select
                label={t('subscription.modals.billing_type')}
                options={[
                  { key: 'token', text: t('subscription.billing_type.token'), value: 'token' },
                  { key: 'request', text: t('subscription.billing_type.request'), value: 'request' },
                ]}
                value={addForm.billing_type}
                onChange={(e, { value }) => setAddForm({ ...addForm, billing_type: value })}
              />
              <Form.Input
                label={t('subscription.modals.duration_days')}
                type='number'
                placeholder={t('subscription.modals.duration_days_placeholder')}
                value={addForm.duration_days}
                onChange={(e, { value }) => setAddForm({ ...addForm, duration_days: value })}
              />
              <Form.TextArea
                label={t('subscription.modals.notes')}
                value={addForm.notes}
                onChange={(e, { value }) => setAddForm({ ...addForm, notes: value })}
              />
            </Form>
          </Modal.Content>
          <Modal.Actions>
            <Button onClick={() => setShowAddModal(false)}>{t('subscription.buttons.cancel')}</Button>
            <Button positive onClick={addSubscription}>{t('subscription.buttons.confirm')}</Button>
          </Modal.Actions>
        </Modal>
      )}

      {/* Renew Modal */}
      {isAdminUser && (
        <Modal open={showRenewModal} onClose={() => setShowRenewModal(false)}>
          <Modal.Header>{t('subscription.modals.renew_title')}</Modal.Header>
          <Modal.Content>
            <Form>
              <Form.Input
                label={t('subscription.modals.new_end_time')}
                type='number'
                value={renewForm.end_time}
                onChange={(e, { value }) => setRenewForm({ ...renewForm, end_time: parseInt(value) })}
              />
              <Form.TextArea
                label={t('subscription.modals.notes')}
                value={renewForm.notes}
                onChange={(e, { value }) => setRenewForm({ ...renewForm, notes: value })}
              />
            </Form>
          </Modal.Content>
          <Modal.Actions>
            <Button onClick={() => setShowRenewModal(false)}>{t('subscription.buttons.cancel')}</Button>
            <Button positive onClick={renewSubscription}>{t('subscription.buttons.confirm')}</Button>
          </Modal.Actions>
        </Modal>
      )}

      {/* Usage Modal */}
      <Modal open={showUsageModal} onClose={() => setShowUsageModal(false)} size='large'>
        <Modal.Header>{t('subscription.modals.usage_title')}</Modal.Header>
        <Modal.Content>
          {usageData && usageData.subscription && (
            <div>
              <p><strong>{t('subscription.table.plan')}:</strong> {usageData.subscription.plan ? usageData.subscription.plan.name : '-'}</p>
              <p><strong>{t('subscription.table.billing_type')}:</strong> {renderBillingType(usageData.billing_type || usageData.subscription.billing_type, t)}</p>
              <p><strong>{t('subscription.table.start_time')}:</strong> {timestamp2string(usageData.subscription.start_time)}</p>
              <p><strong>{t('subscription.table.end_time')}:</strong> {timestamp2string(usageData.subscription.end_time)}</p>
              {renderUsageProgress()}
            </div>
          )}
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={() => setShowUsageModal(false)}>{t('subscription.buttons.close')}</Button>
        </Modal.Actions>
      </Modal>
    </>
  );
};

export default SubscriptionsTable;