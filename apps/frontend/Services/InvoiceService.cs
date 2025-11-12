using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using frontend.Models;
using frontend.Pages;

namespace frontend.Services
{
    public class InvoiceService
    {
        private readonly List<Invoice> _mockInvoices = new();

        public InvoiceService()
        {
            // Seed some mock data
            _mockInvoices.Add(new Invoice
            {
                ID = Guid.NewGuid(),
                InvoiceNumber = "INV-1001",
                CustomerName = "John Doe",
                Amount = 5000,
                Status = "Pending",
                CreatedAt = DateTime.UtcNow.AddDays(-2),
                Items = new List<InvoiceItem>
                {
                    new InvoiceItem { ItemName = "Tomatoes", UnitPrice = 100, Quantity = 20 },
                    new InvoiceItem { ItemName = "Potatoes", UnitPrice = 50, Quantity = 40 }
                }
            });

            _mockInvoices.Add(new Invoice
            {
                ID = Guid.NewGuid(),
                InvoiceNumber = "INV-1002",
                CustomerName = "Mary Jane",
                Amount = 3000,
                Status = "Paid",
                CreatedAt = DateTime.UtcNow.AddDays(-1),
                Items = new List<InvoiceItem>
                {
                    new InvoiceItem { ItemName = "Onions", UnitPrice = 100, Quantity = 20 },
                }
            });
        }

        public Task<List<Invoice>> GetAllAsync()
        {
            return Task.FromResult(_mockInvoices);
        }

        public Task<Invoice?> GetByIdAsync(Guid id)
        {
            var invoice = _mockInvoices.FirstOrDefault(i => i.ID == id);
            return Task.FromResult(invoice);
        }

        public Task<InvoiceSummary> GetSummaryAsync()
        {
            var summary = new InvoiceSummary
            {
                Total = _mockInvoices.Count,
                Pending = _mockInvoices.Count(i => i.Status == "Pending"),
                Paid = _mockInvoices.Count(i => i.Status == "Paid"),
                Overdue = _mockInvoices.Count(i => i.Status == "Overdue"),
                TotalAmount = _mockInvoices.Sum(i => i.Amount)
            };

            return Task.FromResult(summary);
        }

        public Task MarkAsPaidAsync(Guid id)
        {
            var invoice = _mockInvoices.FirstOrDefault(i => i.ID == id);
            if (invoice != null)
            {
                invoice.Status = "Paid";
            }
            return Task.CompletedTask;
        }

        public Task<Invoice> CreateInvoiceAsync(Invoice invoice)
        {
            invoice.ID = Guid.NewGuid();
            invoice.Status = "Pending";
            invoice.CreatedAt = DateTime.UtcNow;
            _mockInvoices.Add(invoice);
            return Task.FromResult(invoice);
        }
    }
}