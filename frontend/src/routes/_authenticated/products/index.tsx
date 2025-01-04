import { useState, useEffect } from 'react';
import { Link, createFileRoute } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import {
  IconPlus,
  IconSearch,
  IconFilter,
  IconArrowUp,
  IconArrowDown,
  IconArrowsUpDown,
} from '@tabler/icons-react';
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { Slider } from '@/components/ui/slider';
import { productsApi } from '@/lib/api/products';
import { Skeleton } from '@/components/ui/skeleton';
import { Pagination } from '@/components/ui/pagination';

export const Route = createFileRoute('/_authenticated/products/')({
  component: ProductsPage,
});

const PER_PAGE = 10;

const categories = [
  { value: 'coffee', label: 'コーヒー豆' },
  { value: 'tea', label: '茶葉' },
  { value: 'equipment', label: '器具・設備' },
  { value: 'accessories', label: 'アクセサリー' },
];

const statuses = [
  { value: '販売中', label: '販売中' },
  { value: '在庫切れ', label: '在庫切れ' },
  { value: '入荷待ち', label: '入荷待ち' },
  { value: '在庫少', label: '在庫少' },
];

function ProductsPage() {
  const [searchQuery, setSearchQuery] = useState('');
  const [debouncedQuery, setDebouncedQuery] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [selectedStatus, setSelectedStatus] = useState<string>('');
  const [priceRange, setPriceRange] = useState<[number, number]>([0, 100000]);
  const [sortBy, setSortBy] = useState<string>('');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  // 検索クエリのデバウンス処理
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(searchQuery);
      setCurrentPage(1); // 検索時はページを1に戻す
    }, 500);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  const { data, isLoading, error } = useQuery({
    queryKey: [
      'products',
      {
        search: debouncedQuery,
        page: currentPage,
        category: selectedCategory,
        status: selectedStatus,
        minPrice: priceRange[0],
        maxPrice: priceRange[1],
        sortBy,
        sortOrder,
      },
    ],
    queryFn: () =>
      productsApi.getProducts({
        search: debouncedQuery,
        page: currentPage,
        perPage: PER_PAGE,
        category: selectedCategory,
        status: selectedStatus,
        minPrice: priceRange[0],
        maxPrice: priceRange[1],
        sortBy,
        sortOrder,
      }),
  });

  const totalPages = data ? Math.ceil(data.total / PER_PAGE) : 0;

  const resetFilters = () => {
    setSelectedCategory('');
    setSelectedStatus('');
    setPriceRange([0, 100000]);
  };

  const handleSort = (field: string) => {
    if (sortBy === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  const getSortIcon = (field: string) => {
    if (sortBy !== field) {
      return <IconArrowsUpDown className="ml-2 h-4 w-4" />;
    }
    return sortOrder === 'asc' ? (
      <IconArrowUp className="ml-2 h-4 w-4" />
    ) : (
      <IconArrowDown className="ml-2 h-4 w-4" />
    );
  };

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
        <div className="flex items-center gap-4">
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

          <Sheet>
            <SheetTrigger asChild>
              <Button variant="outline">
                <IconFilter className="mr-2 h-4 w-4" />
                フィルター
              </Button>
            </SheetTrigger>
            <SheetContent>
              <SheetHeader>
                <SheetTitle>フィルター</SheetTitle>
                <SheetDescription>
                  商品の絞り込み条件を設定します。
                </SheetDescription>
              </SheetHeader>
              <div className="mt-6 space-y-6">
                <div className="space-y-2">
                  <label className="text-sm font-medium">カテゴリー</label>
                  <Select
                    value={selectedCategory}
                    onValueChange={setSelectedCategory}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="カテゴリーを選択" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">すべて</SelectItem>
                      {categories.map((category) => (
                        <SelectItem key={category.value} value={category.value}>
                          {category.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">ステータス</label>
                  <Select value={selectedStatus} onValueChange={setSelectedStatus}>
                    <SelectTrigger>
                      <SelectValue placeholder="ステータスを選択" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">すべて</SelectItem>
                      {statuses.map((status) => (
                        <SelectItem key={status.value} value={status.value}>
                          {status.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium">価格帯</label>
                  <div className="px-2">
                    <Slider
                      value={priceRange}
                      onValueChange={setPriceRange as (value: number[]) => void}
                      max={100000}
                      step={1000}
                      className="my-4"
                    />
                    <div className="flex items-center justify-between text-sm">
                      <span>¥{priceRange[0].toLocaleString()}</span>
                      <span>¥{priceRange[1].toLocaleString()}</span>
                    </div>
                  </div>
                </div>

                <Button
                  variant="outline"
                  className="w-full"
                  onClick={resetFilters}
                >
                  フィルターをリセット
                </Button>
              </div>
            </SheetContent>
          </Sheet>
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
              <TableHead
                className="cursor-pointer"
                onClick={() => handleSort('name')}
              >
                <div className="flex items-center">
                  商品名
                  {getSortIcon('name')}
                </div>
              </TableHead>
              <TableHead
                className="cursor-pointer"
                onClick={() => handleSort('sku')}
              >
                <div className="flex items-center">
                  SKU
                  {getSortIcon('sku')}
                </div>
              </TableHead>
              <TableHead
                className="cursor-pointer text-right"
                onClick={() => handleSort('price')}
              >
                <div className="flex items-center justify-end">
                  価格
                  {getSortIcon('price')}
                </div>
              </TableHead>
              <TableHead
                className="cursor-pointer text-right"
                onClick={() => handleSort('stock')}
              >
                <div className="flex items-center justify-end">
                  在庫数
                  {getSortIcon('stock')}
                </div>
              </TableHead>
              <TableHead
                className="cursor-pointer"
                onClick={() => handleSort('status')}
              >
                <div className="flex items-center">
                  ステータス
                  {getSortIcon('status')}
                </div>
              </TableHead>
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