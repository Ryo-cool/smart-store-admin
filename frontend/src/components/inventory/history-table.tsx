import { InventoryHistory } from '@/lib/api/inventory';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface InventoryHistoryTableProps {
  histories: InventoryHistory[];
}

export function InventoryHistoryTable({ histories }: InventoryHistoryTableProps) {
  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>日時</TableHead>
            <TableHead>種別</TableHead>
            <TableHead className="text-right">数量</TableHead>
            <TableHead>理由</TableHead>
            <TableHead>備考</TableHead>
            <TableHead>担当者</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {histories.map((history) => (
            <TableRow key={history.id}>
              <TableCell>
                {new Date(history.createdAt).toLocaleString('ja-JP')}
              </TableCell>
              <TableCell>
                <span
                  className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                    history.type === '入庫'
                      ? 'bg-green-100 text-green-800'
                      : history.type === '出庫'
                      ? 'bg-red-100 text-red-800'
                      : 'bg-blue-100 text-blue-800'
                  }`}
                >
                  {history.type}
                </span>
              </TableCell>
              <TableCell className="text-right">
                {history.type === '出庫' ? '-' : '+'}
                {history.quantity}
              </TableCell>
              <TableCell>{history.reason}</TableCell>
              <TableCell>{history.note}</TableCell>
              <TableCell>{history.createdBy}</TableCell>
            </TableRow>
          ))}
          {histories.length === 0 && (
            <TableRow>
              <TableCell colSpan={6} className="text-center text-gray-500">
                在庫履歴がありません
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
} 