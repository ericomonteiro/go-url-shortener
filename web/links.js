document.addEventListener('DOMContentLoaded', () => {
    const linksTableBody = document.getElementById('linksTableBody');
    const shortUrlFilter = document.getElementById('shortUrlFilter');
    const destinyUrlFilter = document.getElementById('destinyUrlFilter');
    const loadingElement = document.getElementById('loading');
    const errorElement = document.getElementById('error');
    
    // State variables
    let links = [];
    let sortField = 'shortUrl';
    let sortDirection = 'asc';
    
    // Fetch all links
    async function fetchLinks() {
        try {
            loadingElement.style.display = 'block';
            errorElement.style.display = 'none';
            
            const response = await fetch('/v1/links');
            
            if (!response.ok) {
                throw new Error('Failed to fetch links');
            }
            
            const data = await response.json();
            links = data.links;
            
            renderLinks();
        } catch (error) {
            errorElement.textContent = `Error: ${error.message}`;
            errorElement.style.display = 'block';
        } finally {
            loadingElement.style.display = 'none';
        }
    }
    
    // Render links with current sorting and filtering
    function renderLinks() {
        // Clear table
        linksTableBody.innerHTML = '';
        
        // Apply filters
        const shortUrlFilterValue = shortUrlFilter.value.toLowerCase();
        const destinyUrlFilterValue = destinyUrlFilter.value.toLowerCase();
        
        const filteredLinks = links.filter(link => {
            const shortUrlMatch = !shortUrlFilterValue || link.short_url.toLowerCase().includes(shortUrlFilterValue);
            const destinyUrlMatch = !destinyUrlFilterValue || link.destiny_url.toLowerCase().includes(destinyUrlFilterValue);
            return shortUrlMatch && destinyUrlMatch;
        });
        
        // Apply sorting
        const sortedLinks = [...filteredLinks].sort((a, b) => {
            let valueA, valueB;
            
            switch (sortField) {
                case 'shortUrl':
                    valueA = a.short_url;
                    valueB = b.short_url;
                    break;
                case 'destinyUrl':
                    valueA = a.destiny_url;
                    valueB = b.destiny_url;
                    break;
                case 'clicks':
                    valueA = a.clicks;
                    valueB = b.clicks;
                    // For numeric comparison
                    return sortDirection === 'asc' ? valueA - valueB : valueB - valueA;
                default:
                    valueA = a.short_url;
                    valueB = b.short_url;
            }
            
            // For string comparison
            if (sortField !== 'clicks') {
                if (sortDirection === 'asc') {
                    return valueA.localeCompare(valueB);
                } else {
                    return valueB.localeCompare(valueA);
                }
            }
        });
        
        // Render sorted and filtered links
        sortedLinks.forEach(link => {
            const row = document.createElement('tr');
            
            // Short URL column with link
            const shortUrlCell = document.createElement('td');
            const shortUrlLink = document.createElement('a');
            shortUrlLink.href = link.short_url;
            shortUrlLink.textContent = link.short_url;
            shortUrlLink.target = '_blank';
            shortUrlCell.appendChild(shortUrlLink);
            
            // Destiny URL column with truncated text
            const destinyUrlCell = document.createElement('td');
            destinyUrlCell.className = 'destiny-url';
            destinyUrlCell.textContent = link.destiny_url;
            destinyUrlCell.title = link.destiny_url; // Show full URL on hover
            
            // Clicks column
            const clicksCell = document.createElement('td');
            clicksCell.textContent = link.clicks;
            
            // Add cells to row
            row.appendChild(shortUrlCell);
            row.appendChild(destinyUrlCell);
            row.appendChild(clicksCell);
            
            // Add row to table
            linksTableBody.appendChild(row);
        });
        
        // Update sort icons
        updateSortIcons();
    }
    
    // Update sort icons based on current sort state
    function updateSortIcons() {
        document.querySelectorAll('th.sortable').forEach(th => {
            const field = th.getAttribute('data-sort');
            const icon = th.querySelector('.sort-icon');
            
            if (field === sortField) {
                icon.textContent = sortDirection === 'asc' ? '↑' : '↓';
                th.classList.add('active-sort');
            } else {
                icon.textContent = '↕';
                th.classList.remove('active-sort');
            }
        });
    }
    
    // Add event listeners for sorting
    document.querySelectorAll('th.sortable').forEach(th => {
        th.addEventListener('click', () => {
            const field = th.getAttribute('data-sort');
            
            // If clicking the same field, toggle direction
            if (field === sortField) {
                sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
            } else {
                // New field, default to ascending
                sortField = field;
                sortDirection = 'asc';
            }
            
            renderLinks();
        });
    });
    
    // Add event listeners for filtering
    shortUrlFilter.addEventListener('input', renderLinks);
    destinyUrlFilter.addEventListener('input', renderLinks);
    
    // Initial fetch
    fetchLinks();
});
