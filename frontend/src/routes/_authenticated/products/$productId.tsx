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
import { productsApi } from '@/lib/api/products';

export const Route = createFileRoute('/_authenticated/products/$productId')({
  component: ProductDetailPage,
});

function ProductDetailPage() {
  const { productId } = useParams({ from: Route.fullPath });

  const { data: product, isLoading, error } = useQuery({
    queryKey: ['product', productId],
    queryFn: () => productsApi.getProduct(productId),
  });

  if (isLoading) {
    return <Skeleton className="h-48 w-full" />;
  }

  if (error) {
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
          <Link to="..">
            <Button variant="ghost" size="icon">
              <IconArrowLeft className="h-4 w-4" />
            </Button>
          </Link>
          <div>
            <h1 className="text-2xl font-bold">{product.name}</h1>
            <p className="text-sm text-gray-500">商品の詳細情報</p>
          </div>
        </div>
        <Link to="edit">
          <Button>
            <IconEdit className="mr-2 h-4 w-4" />
            編集
          </Button>
        </Link>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>基本情報</CardTitle>
            <CardDescription>商品の基本的な情報</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
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
    </div>
  );
} 