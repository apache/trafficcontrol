import { Component, OnInit } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { Chart } from 'chart.js';
import { DataPoint, DataSet, DeliveryService } from "../../models";
@Component({
  selector: 'tp-charts',
  templateUrl: './charts.component.html',
  styleUrls: ['./charts.component.scss']
})
export class ChartsComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
    //dummy data change later
    let dataset = [1, 5, 8, 20, 40, 10]
    let min = Math.min(...dataset);
    let max = Math.max(...dataset);
    let average = dataset.reduce((a, b) => a + b, 0) / dataset.length;
    new Chart("minMaxAverageChart", {
      type: 'bar',
      data: {
        labels: ['MIN', 'MAX', 'AVERAGE'],
        datasets: [{
          data: [min, max, average],
          backgroundColor: [
            'rgba(255, 99, 132, 0.2)',
            'rgba(255, 159, 64, 0.2)',
            'rgba(201, 203, 207, 0.2)',
          ],
          borderColor: [
            'rgb(255, 99, 132)',
            'rgb(255, 159, 64)',
            'rgb(201, 203, 207)',
          ],
          borderWidth: 1
        }],
      },
      options: {
        responsive: true,
        legend: {
          display: false,
        },
        scales: {
          yAxes: [{
            ticks: { min: 0, max: 100 },
            gridLines: { display: false }
          }]
        },
      }
    })
    new Chart("percentileChart", {
      type: 'line',
      

    })
  }
}
