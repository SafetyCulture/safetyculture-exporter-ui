import { describe, it, expect } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import DataTableTest from './data-table-test-harness.svelte';

describe('DataTable', () => {
  it('renders all rows', () => {
    render(DataTableTest);
    expect(screen.getByText('Alpha')).toBeInTheDocument();
    expect(screen.getByText('Bravo')).toBeInTheDocument();
    expect(screen.getByText('Charlie')).toBeInTheDocument();
  });

  it('renders column headers', () => {
    render(DataTableTest);
    expect(screen.getByText('Name')).toBeInTheDocument();
    expect(screen.getByText('Date')).toBeInTheDocument();
  });

  it('sorts ascending then descending on header click', async () => {
    render(DataTableTest);

    const sortButton = screen.getByRole('button', { name: /date/i });

    // First click — ascending
    await fireEvent.click(sortButton);
    let rows = screen.getAllByRole('row');
    // row 0 is the header, data rows start at 1
    expect(rows[1]).toHaveTextContent('Charlie');
    expect(rows[2]).toHaveTextContent('Alpha');
    expect(rows[3]).toHaveTextContent('Bravo');

    // Second click — descending
    await fireEvent.click(sortButton);
    rows = screen.getAllByRole('row');
    expect(rows[1]).toHaveTextContent('Bravo');
    expect(rows[2]).toHaveTextContent('Alpha');
    expect(rows[3]).toHaveTextContent('Charlie');
  });

  it('shows "No results." when data is empty', () => {
    render(DataTableTest, { props: { empty: true } });
    expect(screen.getByText('No results.')).toBeInTheDocument();
  });
});
