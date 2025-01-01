import { useState } from 'react'
import { Link, createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { IconEdit, IconEye } from '@tabler/icons-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Skeleton } from '@/components/ui/skeleton'
import { deliveriesApi } from '@/lib/api/deliveries'

interface Delivery {
  id: string
  deliveryType: string
  address: string
  estimatedDeliveryTime: string
  status: string
}

export const Route = createFileRoute('/_authenticated/deliveries/')({
  component: DeliveriesPage,
})

function DeliveriesPage() {
  const [search, setSearch] = useState('')
  const [status, setStatus] = useState<string | undefined>()

  const { data, isLoading } = useQuery({
    queryKey: ['deliveries', { search, status }],
    queryFn: () => deliveriesApi.getDeliveries({ search, status }),
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">配送管理</h1>
      </div>

      <div className="flex items-center gap-4">
        <Input
          placeholder="配送IDまたは配送先で検索"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="max-w-xs"
        />
        <Select value={status} onValueChange={setStatus}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="ステータスで絞り込み" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">すべて</SelectItem>
            <SelectItem value="配送準備中">配送準備中</SelectItem>
            <SelectItem value="配送中">配送中</SelectItem>
            <SelectItem value="配送完了">配送完了</SelectItem>
            <SelectItem value="配送失敗">配送失敗</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {isLoading ? (
        <Skeleton className="h-96 w-full" />
      ) : (
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>配送ID</TableHead>
                <TableHead>配送方法</TableHead>
                <TableHead>配送先</TableHead>
                <TableHead>配送予定日時</TableHead>
                <TableHead>ステータス</TableHead>
                <TableHead className="w-[100px]">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data?.deliveries.map((delivery: Delivery) => (
                <TableRow key={delivery.id}>
                  <TableCell>{delivery.id}</TableCell>
                  <TableCell>{delivery.deliveryType}</TableCell>
                  <TableCell>{delivery.address}</TableCell>
                  <TableCell>
                    {new Date(delivery.estimatedDeliveryTime).toLocaleString(
                      'ja-JP',
                    )}
                  </TableCell>
                  <TableCell>
                    <span
                      className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
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
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <Link
                        to="$deliveryId"
                        params={{ deliveryId: delivery.id }}
                      >
                        <Button variant="ghost" size="icon">
                          <IconEye className="h-4 w-4" />
                        </Button>
                      </Link>
                      <Link
                        to="$deliveryId/edit"
                        params={{ deliveryId: delivery.id }}
                      >
                        <Button variant="ghost" size="icon">
                          <IconEdit className="h-4 w-4" />
                        </Button>
                      </Link>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}
