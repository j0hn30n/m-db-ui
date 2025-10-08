// 全局工具函数
function showError(message) {
    document.getElementById('errorMessage').textContent = message;
    const modal = new bootstrap.Modal(document.getElementById('errorModal'));
    modal.show();
}

function showSuccess(message) {
    document.getElementById('successMessage').textContent = message;
    const toast = new bootstrap.Toast(document.getElementById('successToast'));
    toast.show();
}

// 格式化JSON
function formatJSON(obj) {
    return JSON.stringify(obj, null, 2);
}

// 格式化字节大小
function formatBytes(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 复制到剪贴板
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showSuccess('已复制到剪贴板');
    }).catch(err => {
        showError('复制失败: ' + err.message);
    });
}

// 确认对话框
function confirmAction(message, callback) {
    if (confirm(message)) {
        callback();
    }
}

// 显示统计信息
function showStats() {
    fetch('/api/v1/stats')
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            showError(data.error);
            return;
        }

        const statsHtml = `
            <div class="row">
                <div class="col-md-6">
                    <div class="card">
                        <div class="card-header">
                            <h6 class="mb-0">服务器信息</h6>
                        </div>
                        <div class="card-body">
                            <p><strong>版本:</strong> ${data.version}</p>
                            <p><strong>运行时间:</strong> ${Math.floor(data.uptime / 3600)} 小时</p>
                            <p><strong>数据库数量:</strong> ${data.databaseCount}</p>
                        </div>
                    </div>
                </div>
                <div class="col-md-6">
                    <div class="card">
                        <div class="card-header">
                            <h6 class="mb-0">连接信息</h6>
                        </div>
                        <div class="card-body">
                            <p><strong>当前连接:</strong> ${data.connections.current}</p>
                            <p><strong>可用连接:</strong> ${data.connections.available}</p>
                            <p><strong>总创建连接:</strong> ${data.connections.totalCreated}</p>
                        </div>
                    </div>
                </div>
                <div class="col-md-12 mt-3">
                    <div class="card">
                        <div class="card-header">
                            <h6 class="mb-0">内存使用</h6>
                        </div>
                        <div class="card-body">
                            <p><strong>常驻内存:</strong> ${formatBytes(data.memory.resident * 1024 * 1024)}</p>
                            <p><strong>虚拟内存:</strong> ${formatBytes(data.memory.virtual * 1024 * 1024)}</p>
                            ${data.memory.mapped ? `<p><strong>映射内存:</strong> ${formatBytes(data.memory.mapped * 1024 * 1024)}</p>` : ''}
                        </div>
                    </div>
                </div>
            </div>
        `;

        // 创建模态框显示统计信息
        const modal = document.createElement('div');
        modal.className = 'modal fade';
        modal.innerHTML = `
            <div class="modal-dialog modal-lg">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">
                            <i class="fas fa-chart-bar me-2"></i>服务器统计信息
                        </h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        ${statsHtml}
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">关闭</button>
                    </div>
                </div>
            </div>
        `;
        document.body.appendChild(modal);
        const bsModal = new bootstrap.Modal(modal);
        bsModal.show();

        // 模态框关闭时移除元素
        modal.addEventListener('hidden.bs.modal', () => {
            document.body.removeChild(modal);
        });
    })
    .catch(error => showError('获取统计信息失败: ' + error.message));
}

// 表格工具函数
function sortTable(table, column, asc = true) {
    const tbody = table.querySelector('tbody');
    const rows = Array.from(tbody.querySelectorAll('tr'));

    rows.sort((a, b) => {
        const aValue = a.cells[column].textContent.trim();
        const bValue = b.cells[column].textContent.trim();

        if (asc) {
            return aValue.localeCompare(bValue);
        } else {
            return bValue.localeCompare(aValue);
        }
    });

    rows.forEach(row => tbody.appendChild(row));
}

// 搜索功能
function searchTable(table, searchTerm) {
    const tbody = table.querySelector('tbody');
    const rows = tbody.querySelectorAll('tr');

    rows.forEach(row => {
        const text = row.textContent.toLowerCase();
        const matches = text.includes(searchTerm.toLowerCase());
        row.style.display = matches ? '' : 'none';
    });
}

// 导出功能
function exportToJSON(data, filename) {
    const blob = new Blob([formatJSON(data)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename || 'export.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

// 导入功能
function importFromJSON(callback) {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json';
    input.onchange = function(event) {
        const file = event.target.files[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = function(e) {
            try {
                const data = JSON.parse(e.target.result);
                callback(data);
            } catch (error) {
                showError('无效的JSON文件: ' + error.message);
            }
        };
        reader.readAsText(file);
    };
    input.click();
}

// 页面加载完成后的初始化
document.addEventListener('DOMContentLoaded', function() {
    // 添加表格搜索功能
    const searchInputs = document.querySelectorAll('input[type="search"]');
    searchInputs.forEach(input => {
        input.addEventListener('input', function() {
            const table = this.closest('.card-body').querySelector('table');
            if (table) {
                searchTable(table, this.value);
            }
        });
    });

    // 添加复制按钮
    const codeBlocks = document.querySelectorAll('pre code');
    codeBlocks.forEach(block => {
        const button = document.createElement('button');
        button.className = 'btn btn-sm btn-outline-secondary position-absolute top-0 end-0 m-2';
        button.innerHTML = '<i class="fas fa-copy"></i>';
        button.onclick = () => copyToClipboard(block.textContent);

        const container = block.parentElement;
        container.style.position = 'relative';
        container.appendChild(button);
    });
});

// 错误处理
window.addEventListener('error', function(event) {
    console.error('Global error:', event.error);
    showError('发生了一个错误: ' + event.error.message);
});

// 网络错误处理
window.addEventListener('unhandledrejection', function(event) {
    console.error('Unhandled promise rejection:', event.reason);
    showError('网络请求失败: ' + event.reason);
});