import React from 'react';
import { Card } from 'semantic-ui-react';
import { useTranslation } from 'react-i18next';
import SubscriptionsTable from '../../components/SubscriptionsTable';

const Subscription = () => {
  const { t } = useTranslation();

  return (
    <div className='dashboard-container'>
      <Card fluid className='chart-card'>
        <Card.Content>
          <Card.Header className='header'>{t('subscription.title')}</Card.Header>
          <SubscriptionsTable />
        </Card.Content>
      </Card>
    </div>
  );
};

export default Subscription;