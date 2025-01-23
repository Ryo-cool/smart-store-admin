import { Link, useParams, createFileRoute } from '@tanstack/react-router';
import { IconArrowLeft, IconPackage } from '@tabler/icons-react';
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
import { productsApi } from '@/lib/api/products';
import { inventoryApi } from '@/lib/api/inventory';
import { InventoryHistoryTable } from '@/components/inventory/history-table';
import { InventoryUpdateDialog } from '@/components/inventory/update-dialog';

export const Route = createFileRoute('/_authenticated/products/$productId')({
  component: ProductDetailPage,
});

function ProductDetailPage() {
  const { productId } = useParams({ from: Route.fullPath });

  const { data: product, isLoading: isLoadingProduct, error: productError } = useQuery({
    queryKey: ['product', productId],
    queryFn: () => productsApi.getProduct(productId),
  });

  const { data: inventoryHistory, isLoading: isLoadingHistory } = useQuery({
    queryKey: ['inventory', productId],
    queryFn: () => inventoryApi.getProductHistory(productId),
  });

  if (isLoadingProduct) {
    return <Skeleton className="h-48 w-full" />;
  }

  if (productError) {
    return <div className="text-center text-red-500">商品の取得に失敗しました</div>;
  }

  if (!product) {
    return <div className="text-center">商品が見つかりませんでした</div>;
  }

  const details = [
    { label: 'SKU', value: product.sku },
    { label: 'カテゴリー', value: product.category },
    { label: '重量', value: product.weight },
    { label: 'サイズ', value: product.dimensions },
    { label: '在庫数', value: product.stock },
    { label: '登録日', value: product.createdAt },
    { label: '更新日', value: product.updatedAt },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Link to=".." search={{}}>
            <Button variant="ghost" size="icon">
              <IconArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold">{product.name}</h1>
            <p className="text-sm text-gray-500">商品の詳細情報</p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <InventoryUpdateDialog
            productId={productId}
            productName={product.name}
            trigger={
              <Button variant="outline">
                <IconPackage className="mr-2 h-4 w-4" />
                在庫更新
              </Button>
            }
          />
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>基本情報</CardTitle>
            <CardDescription>商品の基本的な情報</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {product.images && product.images.length > 0 && (
              <div>
                <h3 className="text-sm font-medium text-gray-500">商品画像</h3>
                <div className="mt-2 grid grid-cols-2 gap-4 sm:grid-cols-3">
                  {product.images.map((image, index) => (
                    <div key={index} className="relative aspect-square">
                      <img
                        src={image}
                        alt={`${product.name} - 画像 ${index + 1}`}
                        className="h-full w-full rounded-lg object-cover"
                      />
                    </div>
                  ))}
                </div>
              </div>
            )}
            <div>
              <h3 className="text-sm font-medium text-gray-500">商品名</h3>
              <p className="mt-1">{product.name}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500">説明</h3>
              <p className="mt-1">{product.description}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500">価格</h3>
              <p className="mt-1">¥{product.price.toLocaleString()}</p>
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-500">ステータス</h3>
              <span
                className={`mt-1 inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                  product.status === '在庫少'
                    ? 'bg-yellow-100 text-yellow-800'
                    : product.status === '在庫切れ'
                    ? 'bg-red-100 text-red-800'
                    : product.status === '入荷待ち'
                    ? 'bg-blue-100 text-blue-800'
                    : 'bg-green-100 text-green-800'
                }`}
              >
                {product.status}
              </span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>詳細情報</CardTitle>
            <CardDescription>商品の詳細な情報</CardDescription>
          </CardHeader>
          <CardContent>
            <dl className="divide-y">
              {details.map((detail) => (
                <div
                  key={detail.label}
                  className="flex justify-between py-3 text-sm"
                >
                  <dt className="text-gray-500">{detail.label}</dt>
                  <dd className="font-medium text-gray-900">{detail.value}</dd>
                </div>
              ))}
            </dl>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>在庫履歴</CardTitle>
          <CardDescription>商品の在庫変動履歴</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoadingHistory ? (
            <Skeleton className="h-48 w-full" />
          ) : (
            <InventoryHistoryTable histories={inventoryHistory?.histories || []} />
          )}
        </CardContent>
      </Card>
    </div>
  );
} 