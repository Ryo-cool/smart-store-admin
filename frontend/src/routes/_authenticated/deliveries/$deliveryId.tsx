import { Link, useParams, createFileRoute } from '@tanstack/react-router';
import { IconArrowLeft, IconEdit } from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { deliveriesApi } from '@/lib/api/deliveries';
import { DeliveryStatusDialog } from '@/components/delivery/status-dialog';

interface DeliveryHistoryItem {
  status: string;
  timestamp: string;
  location?: {
    latitude: number;
    longitude: number;
  };
  note?: string;
}

export const Route = createFileRoute('/_authenticated/deliveries/$deliveryId')({
  component: DeliveryDetailPage,
});

function DeliveryDetailPage() {
  const { deliveryId } = useParams({ from: '/_authenticated/deliveries/$deliveryId' });

  const { data: delivery, isLoading: isLoadingDelivery } = useQuery({
    queryKey: ['delivery', deliveryId],
    queryFn: () => deliveriesApi.getDelivery(deliveryId),
  });

  const { data: history, isLoading: isLoadingHistory } = useQuery({
    queryKey: ['delivery-history', deliveryId],
    queryFn: () => deliveriesApi.getDeliveryHistory(deliveryId),
  });

  if (isLoadingDelivery) {
    return <Skeleton className="h-48 w-full" />;
  }

  if (!delivery) {
    return <div className="text-center">配送情報が見つかりませんでした</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link to="..">
            <Button variant="ghost" size="icon">
              <IconArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold">配送詳細</h1>
            <p className="text-sm text-gray-500">配送ID: {delivery.id}</p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <DeliveryStatusDialog
            deliveryId={deliveryId}
            currentStatus={delivery.status}
            trigger={
              <Button variant="outline">
                <IconEdit className="mr-2 h-4 w-4" />
                ステータス更新
              </Button>
            }
          />
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>基本情報</CardTitle>
            <CardDescription>配送の基本的な情報</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <h3 className="text-sm font-medium text-gray-500">配送方法</h3>
              <p className="mt-1">{delivery.deliveryType}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500">配送先</h3>
              <p className="mt-1">{delivery.address}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500">配送予定日時</h3>
              <p className="mt-1">
                {new Date(delivery.estimatedDeliveryTime).toLocaleString('ja-JP')}
              </p>
            </div>
            {delivery.actualDeliveryTime && (
              <div>
                <h3 className="text-sm font-medium text-gray-500">
                  実際の配送日時
                </h3>
                <p className="mt-1">
                  {new Date(delivery.actualDeliveryTime).toLocaleString('ja-JP')}
                </p>
              </div>
            )}
            <div>
              <h3 className="text-sm font-medium text-gray-500">ステータス</h3>
              <span
                className={`mt-1 inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                  delivery.status === '配送中'
                    ? 'bg-blue-100 text-blue-800'
                    : delivery.status === '配送完了'
                    ? 'bg-green-100 text-green-800'
                    : delivery.status === '配送失敗'
                    ? 'bg-red-100 text-red-800'
                    : 'bg-gray-100 text-gray-800'
                }`}
              >
                {delivery.status}
              </span>
            </div>
            {delivery.notes && (
              <div>
                <h3 className="text-sm font-medium text-gray-500">備考</h3>
                <p className="mt-1">{delivery.notes}</p>
              </div>
            )}
          </CardContent>
        </Card>

        {delivery.trackingInfo && (
          <Card>
            <CardHeader>
              <CardTitle>配送状況</CardTitle>
              <CardDescription>現在の配送状況の詳細</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {delivery.trackingInfo.currentLocation && (
                <div>
                  <h3 className="text-sm font-medium text-gray-500">現在地</h3>
                  <p className="mt-1">
                    緯度: {delivery.trackingInfo.currentLocation.latitude}
                    <br />
                    経度: {delivery.trackingInfo.currentLocation.longitude}
                  </p>
                </div>
              )}
              {delivery.trackingInfo.batteryLevel !== undefined && (
                <div>
                  <h3 className="text-sm font-medium text-gray-500">
                    バッテリー残量
                  </h3>
                  <p className="mt-1">{delivery.trackingInfo.batteryLevel}%</p>
                </div>
              )}
              {delivery.trackingInfo.speed !== undefined && (
                <div>
                  <h3 className="text-sm font-medium text-gray-500">
                    現在の速度
                  </h3>
                  <p className="mt-1">{delivery.trackingInfo.speed} km/h</p>
                </div>
              )}
            </CardContent>
          </Card>
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>配送履歴</CardTitle>
          <CardDescription>配送状況の変更履歴</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoadingHistory ? (
            <Skeleton className="h-48 w-full" />
          ) : (
            <div className="space-y-4">
              {history?.history.map((item: DeliveryHistoryItem, index: number) => (
                <div
                  key={index}
                  className="flex items-start gap-4 border-l-2 border-gray-200 pl-4"
                >
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span
                        className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                          item.status === '配送中'
                            ? 'bg-blue-100 text-blue-800'
                            : item.status === '配送完了'
                            ? 'bg-green-100 text-green-800'
                            : item.status === '配送失敗'
                            ? 'bg-red-100 text-red-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {item.status}
                      </span>
                      <span className="text-sm text-gray-500">
                        {new Date(item.timestamp).toLocaleString('ja-JP')}
                      </span>
                    </div>
                    {item.location && (
                      <p className="mt-1 text-sm">
                        位置情報: {item.location.latitude},{' '}
                        {item.location.longitude}
                      </p>
                    )}
                    {item.note && (
                      <p className="mt-1 text-sm text-gray-600">{item.note}</p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
} 