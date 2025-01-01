import { useState, useEffect } from 'react';
import { Link } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import { IconPlus, IconSearch } from '@tabler/icons-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { productsApi } from '@/lib/api/products';
import { Skeleton } from '@/components/ui/skeleton';
import { Pagination } from '@/components/ui/pagination';

const PER_PAGE = 10;

export default function ProductsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [debouncedQuery, setDebouncedQuery] = useState('');
  const [currentPage, setCurrentPage] = useState(1);

  // 検索クエリのデバウンス処理
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(searchQuery);
      setCurrentPage(1); // 検索時はページを1に戻す
    }, 500);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  const { data, isLoading, error } = useQuery({
    queryKey: ['products', { search: debouncedQuery, page: currentPage }],
    queryFn: () =>
      productsApi.getProducts({
        search: debouncedQuery,
        page: currentPage,
        perPage: PER_PAGE,
      }),
  });

  const totalPages = data ? Math.ceil(data.total / PER_PAGE) : 0;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">商品管理</h1>
          <p className="text-sm text-gray-500">商品の一覧と管理</p>
        </div>
        <Link to="/products/new">
          <Button>
            <IconPlus className="mr-2 h-4 w-4" />
            新規商品
          </Button>
        </Link>
      </div>

      <div className="flex items-center justify-between gap-4">
        <div className="relative w-64">
          <IconSearch className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
          <Input
            type="search"
            placeholder="商品を検索..."
            className="pl-9"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        {data && (
          <div className="text-sm text-gray-500">
            全{data.total}件中 {(currentPage - 1) * PER_PAGE + 1}-
            {Math.min(currentPage * PER_PAGE, data.total)}件を表示
          </div>
        )}
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>商品名</TableHead>
              <TableHead>SKU</TableHead>
              <TableHead className="text-right">価格</TableHead>
              <TableHead className="text-right">在庫数</TableHead>
              <TableHead>ステータス</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              // ローディング時のスケルトン表示
              Array.from({ length: PER_PAGE }).map((_, index) => (
                <TableRow key={index}>
                  <TableCell>
                    <Skeleton className="h-4 w-[200px]" />
                  </TableCell>
                  <TableCell>
                    <Skeleton className="h-4 w-[100px]" />
                  </TableCell>
                  <TableCell className="text-right">
                    <Skeleton className="ml-auto h-4 w-[80px]" />
                  </TableCell>
                  <TableCell className="text-right">
                    <Skeleton className="ml-auto h-4 w-[60px]" />
                  </TableCell>
                  <TableCell>
                    <Skeleton className="h-4 w-[80px]" />
                  </TableCell>
                </TableRow>
              ))
            ) : error ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center text-red-500">
                  データの取得に失敗しました。
                </TableCell>
              </TableRow>
            ) : (
              data?.products.map((product) => (
                <TableRow key={product.id}>
                  <TableCell>
                    <Link
                      to="/products/$productId"
                      params={{ productId: product.id }}
                      className="text-blue-600 hover:underline"
                    >
                      {product.name}
                    </Link>
                  </TableCell>
                  <TableCell>{product.sku}</TableCell>
                  <TableCell className="text-right">
                    ¥{product.price.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-right">{product.stock}</TableCell>
                  <TableCell>
                    <span
                      className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
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
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {totalPages > 1 && (
        <div className="mt-4">
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={setCurrentPage}
          />
        </div>
      )}
    </div>
  );
} 