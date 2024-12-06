function initChart() {
    const chart = document.getElementById('chart');
    if (!chart) return;
    
    const verticalLine = document.getElementById('vertical-line'); // Vertical line for the SVG chart
    const tooltipWidth = parseInt(parseFloat(chart.dataset.tooltipWidth), 10); // Width of the tooltip
    const svgWidth = parseInt(chart.viewBox.baseVal.width, 10); // Width of the SVG element
    let debounceTimeout; // Timeout for debouncing mousemove events

    // Resize observer to update the chart scale
    const resizeObserver = new ResizeObserver(entries => {
        for (const entry of entries) {
            if (entry.target === chart) {
                updateChartScale();
            }
        }
    });

    resizeObserver.observe(chart);

    function updateChartScale() {
        const rect = chart.getBoundingClientRect();
        chart.style.height = `${rect.width * 0.5}px`;
    }

    // Convert mouse position to SVG scale
    function getScaledMousePosition(event, svg) {
        const rect = svg.getBoundingClientRect();
        const mouseX = event.clientX - rect.left;
        const scale = svgWidth / rect.width;
        return mouseX * scale;
    }

    function updateVerticalLine(x) {
        verticalLine.setAttribute('x1', x);
        verticalLine.setAttribute('x2', x);
        verticalLine.classList.remove('opacity-0');
    }

    function positionTooltip(pointX, tooltip) {
        const rect = tooltip.querySelector('rect');
        const text = tooltip.querySelector('text');
        
        if (pointX < tooltipWidth / 2) {
            rect.setAttribute('transform', 'translate(0,0)');
            text.setAttribute('x', pointX + tooltipWidth / 2);
        } else if (pointX > svgWidth - tooltipWidth) {
            rect.setAttribute('transform', `translate(-${tooltipWidth},0)`);
            text.setAttribute('x', pointX - tooltipWidth / 2);
        } else {
            rect.setAttribute('transform', `translate(-${tooltipWidth / 2},0)`);
            text.setAttribute('x', pointX);
        }
        
        const valueY = parseFloat(rect.getAttribute('y'));
        if (valueY < 0) {
            rect.setAttribute('y', '5');
            text.setAttribute('y', '20');
        }
    }

    function handleMouseMove(e) {
        if (debounceTimeout) clearTimeout(debounceTimeout);
        
        debounceTimeout = setTimeout(() => {
            const mouseX = getScaledMousePosition(e, this);
            updateVerticalLine(mouseX);
            
            const points = [...document.querySelectorAll('.chart-point')];
            const spacing = this.viewBox.baseVal.width / (points.length - 1);
            const threshold = spacing * 0.5;
            
            points.forEach(point => {
                const pointX = parseFloat(point.dataset.x);
                const tooltip = point.querySelector('.tooltip');
                const circle = point.querySelector('circle');
                
                if (Math.abs(pointX - mouseX) < threshold) {
                    tooltip.classList.remove('opacity-0');
                    circle.classList.remove('opacity-0');
                    positionTooltip(pointX, tooltip);
                } else {
                    tooltip.classList.add('opacity-0');
                    circle.classList.add('opacity-0');
                }
            });
        }, );
    }

    function handleMouseLeave() {
        if (debounceTimeout) clearTimeout(debounceTimeout);
        verticalLine.classList.add('opacity-0');
        document.querySelectorAll('.chart-point .tooltip, .chart-point circle').forEach(el => {
            el.classList.add('opacity-0');
        });
    }

    chart.addEventListener('mousemove', handleMouseMove);
    chart.addEventListener('mouseleave', handleMouseLeave);

    return () => {
        resizeObserver.disconnect();
        chart.removeEventListener('mousemove', handleMouseMove);
        chart.removeEventListener('mouseleave', handleMouseLeave);
    };
}

// Initialize on page load and HTMX swaps
document.addEventListener('htmx:afterSwap', () => {
    initChart();
});

// Initial load
initChart();