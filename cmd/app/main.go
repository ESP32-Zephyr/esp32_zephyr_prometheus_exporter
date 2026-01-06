package main

import (
        "net/http"
        "time"
        "fmt"
        
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"

        "github.com/ESP32-Zephyr/esp32_zephyr_goapi/api"
)

func getMetrics(transport, ipv4 string, destPort uint16, samplingPeriod time.Duration) {
        gaugeAdcChs := []prometheus.Gauge{}

        es32client, err := api.NewEsp32Client(transport, ipv4, destPort)
        if err != nil {
                fmt.Println(err)
        }

        var chs uint32
        for {
                adcChs, err := es32client.AdcChsGet()        
                if err != nil {
                        fmt.Println("Error getting ADC channels:", err)
                        // Wait indefinetly until the device responds
                        time.Sleep(5 * time.Second)
                } else {
                        chs = adcChs.GetAdcChs()
                        fmt.Println("ADC Channels:", chs)                        
                        if chs == 0 {
                                fmt.Println("No ADC channels available, retrying...")
                                time.Sleep(5 * time.Second)
                                continue
                        }
                        break
                }
        }

        for ch := range chs {
                adcGauge := prometheus.NewGauge(
                        prometheus.GaugeOpts{
                        Name: fmt.Sprintf("esp32_ch_%d", ch),
                        Help: fmt.Sprintf("ADC channel %d value.", ch),
                        },
                )
                gaugeAdcChs = append(gaugeAdcChs, adcGauge)
                prometheus.MustRegister(adcGauge)
        }


        for{
                for id, gaugeAdcCh := range gaugeAdcChs {
                        adc_val, err := es32client.AdcChRead(uint32(id))
                        if err != nil {
                                fmt.Println("Error reading ADC channel:", err)
                                continue
                        }
                        gaugeAdcCh.Set(float64(adc_val.GetVal()))
                }
                time.Sleep(samplingPeriod * time.Second)
        }
}

func main() {       
        transport := "tcp"
        hosts := []string{"192.168.0.4"}
        destPort := uint16(4242)
        const samplingPeriod time.Duration = 5 // in seconds
	
        for _, host := range hosts {
                go getMetrics(transport, host, destPort, samplingPeriod)
        }

        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
}
