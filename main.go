package main

import (
    "machine"
    "time"
    //"device/avr" // потребуется для использования вставок с ассемблером 
)

// массивы с объектами выводов
var arrX [8]machine.Pin = [8]machine.Pin{machine.PB0, machine.PB1, machine.PB2, machine.PB3, machine.PB4, machine.PB5, machine.PB6, machine.PB7}
var arrY [8]machine.Pin = [8]machine.Pin{machine.PD0, machine.PD1, machine.PD2, machine.PD3, machine.PD4, machine.PD5, machine.PD6, machine.PD7}

// функция инициализаци выводов
func Init(x [8]machine.Pin, y [8]machine.Pin) {
	for i := 0; i < 8; i++ {
		x[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
		x[i].Low()
		y[i].Configure(machine.PinConfig{Mode: machine.PinOutput})
		y[i].High()
	}
}

// функция, преобразующая массив в изображение на матрице
func OutMatrix(array []uint8, x int8) {
	for i:=0; i < 8; i++ {
		for j:=0; j < 8; j++ {
			arrX[uint(j)].Set(GetBit(array[int8(i)+x], uint(j)))
		}	
		arrY[i].Low()
		time.Sleep(time.Millisecond * 1)
		arrY[i].High()
	}
}

// функция извлещения бита из байта
func GetBit(n uint8, pos uint) bool {
	mask := uint8(1<<pos)
	result := n & mask
	return result != 0
}


func main() {
	// Инициализируем выводы
	Init(arrX, arrY)
	
	// Объявляем массив с бегущей строкой
	array := []uint8{0x81, 0x42, 0x24, 0x18, 0x18, 0x24, 0x42, 0x81,
			0xFF, 0xC3, 0xA5, 0x99, 0x99, 0xA5, 0xC3, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0x18, 0x18, 0x18, 0xFF, 0xFF, 0x18, 0x18, 0x18}
			
	// Счетчик. С каждым инкрементом смещает каретку бегущей строки
	var counter int8 = 0
	
	// Переменная хранения скорости строки
	var del int8 = 0

	// Переменная хранения направления
	var direction int8 = 1
	
	// Инициализируем кнопку для перекулючения скоростей
	button1 := machine.PC4
	button1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	
	
	for {
		for i := 0; i < 2 + int(del); i++{
			// Передаем массив и каретку(смещение по массиву)
			OutMatrix(array, counter)
			
			// Проверяем нажатие кнопки
			// Для избежания помех, после нажатия надо дождаться момента отпускания
			if !button1.Get() {
				time.Sleep(time.Millisecond * 10)
				for !button1.Get() {
					// Ждем когда кнопку отпустят...
				}
				if del == 0 {
					del = 10
				} else {
					del = 0
				}
			}
		}

		// Смещаем каретку
		counter = counter+direction

		// Если каретка достигла конца массива или начала, то меняем направление 
		if counter == int8(len(array))-8 || counter == 0 {
			direction = -direction
		}
	}
}







